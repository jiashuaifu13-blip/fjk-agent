package conversation

import (
	"context"
	"fmt"
	"strings"

	"github.com/fjk-dev/fjk-agent/core/event"
	"github.com/fjk-dev/fjk-agent/core/model"
)

type Service struct {
	manager  *Manager
	provider model.Provider
	events   *event.Bus
}

func NewService(
	manager *Manager,
	provider model.Provider,
	events *event.Bus,
) *Service {
	return &Service{
		manager:  manager,
		provider: provider,
		events:   events,
	}
}

func (s *Service) Chat(
	ctx context.Context,
	conversationID string,
	userContent string,
) (*Message, error) {
	userMessage := s.manager.AddMessage(
		ctx,
		conversationID,
		RoleUser,
		userContent,
	)

	s.events.Publish(
		ctx,
		event.New(
			event.TypeConversationMessageCreated,
			event.Metadata{
				SessionID: conversationID,
			},
			map[string]any{
				"message_id": userMessage.ID,
				"role":       userMessage.Role,
				"content":    userMessage.Content,
			},
		),
	)

	history := s.manager.Messages(
		ctx,
		conversationID,
	)

	requestMessages := make(
		[]model.Message,
		0,
		len(history),
	)

	for _, message := range history {
		requestMessages = append(
			requestMessages,
			model.Message{
				Role:    model.Role(message.Role),
				Content: message.Content,
			},
		)
	}

	var builder strings.Builder

	err := s.provider.Stream(
		ctx,
		model.Request{
			Messages: requestMessages,
		},
		func(
			ctx context.Context,
			chunk string,
		) error {
			builder.WriteString(chunk)

			s.events.Publish(
				ctx,
				event.New(
					event.TypeConversationMessageDelta,
					event.Metadata{
						SessionID: conversationID,
					},
					map[string]any{
						"delta": chunk,
					},
				),
			)

			return nil
		},
	)

	if err != nil {
		s.events.Publish(
			ctx,
			event.New(
				event.TypeConversationFailed,
				event.Metadata{
					SessionID: conversationID,
				},
				map[string]any{
					"error": err.Error(),
				},
			),
		)

		return nil, fmt.Errorf(
			"conversation stream failed: %w",
			err,
		)
	}

	assistantMessage := s.manager.AddMessage(
		ctx,
		conversationID,
		RoleAssistant,
		builder.String(),
	)

	s.events.Publish(
		ctx,
		event.New(
			event.TypeConversationMessageCompleted,
			event.Metadata{
				SessionID: conversationID,
			},
			map[string]any{
				"message_id": assistantMessage.ID,
				"role":       assistantMessage.Role,
				"content":    assistantMessage.Content,
			},
		),
	)

	return assistantMessage, nil
}
