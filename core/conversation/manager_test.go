package conversation

import (
	"context"
	"testing"
)

func TestConversationMessages(t *testing.T) {
	ctx := context.Background()

	manager := NewManager()

	conversation :=
		manager.Create(
			ctx,
			"Test Conversation",
		)

	manager.AddMessage(
		ctx,
		conversation.ID,
		RoleUser,
		"hello",
	)

	manager.AddMessage(
		ctx,
		conversation.ID,
		RoleAssistant,
		"hi",
	)

	messages := manager.Messages(
		ctx,
		conversation.ID,
	)

	if len(messages) != 2 {
		t.Fatalf(
			"expected 2 messages, got %d",
			len(messages),
		)
	}
}
