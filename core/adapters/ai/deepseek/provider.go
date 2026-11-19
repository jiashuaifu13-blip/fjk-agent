package deepseek

import (
	"context"
	"fmt"

	einodeepseek "github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/schema"

	fjkmodel "github.com/fjk-dev/fjk-agent/core/model"
)

type Config struct {
	APIKey  string
	BaseURL string
	Model   string
}

type Provider struct {
	chatModel *einodeepseek.ChatModel
}

func New(
	ctx context.Context,
	cfg Config,
) (*Provider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf(
			"deepseek API key is required",
		)
	}

	if cfg.Model == "" {
		return nil, fmt.Errorf(
			"deepseek model is required",
		)
	}

	chatModel, err :=
		einodeepseek.NewChatModel(
			ctx,
			&einodeepseek.ChatModelConfig{
				APIKey:  cfg.APIKey,
				BaseURL: cfg.BaseURL,
				Model:   cfg.Model,
			},
		)

	if err != nil {
		return nil, fmt.Errorf(
			"create deepseek chat model: %w",
			err,
		)
	}

	return &Provider{
		chatModel: chatModel,
	}, nil
}

func (p *Provider) Name() string {
	return "deepseek"
}

func (p *Provider) Generate(
	ctx context.Context,
	request fjkmodel.Request,
) (*fjkmodel.Response, error) {
	messages := make(
		[]*schema.Message,
		0,
		len(request.Messages),
	)

	for _, message := range request.Messages {
		var role schema.RoleType

		switch message.Role {
		case fjkmodel.RoleSystem:
			role = schema.System

		case fjkmodel.RoleUser:
			role = schema.User

		case fjkmodel.RoleAssistant:
			role = schema.Assistant

		case fjkmodel.RoleTool:
			role = schema.Tool

		default:
			return nil, fmt.Errorf(
				"unsupported message role: %s",
				message.Role,
			)
		}

		messages = append(
			messages,
			&schema.Message{
				Role:    role,
				Content: message.Content,
			},
		)
	}

	response, err :=
		p.chatModel.Generate(
			ctx,
			messages,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"deepseek generate: %w",
			err,
		)
	}

	return &fjkmodel.Response{
		Content: response.Content,
	}, nil
}
