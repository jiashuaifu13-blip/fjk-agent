package model

import "context"

type StreamHandler func(
	ctx context.Context,
	chunk string,
) error

type Provider interface {
	Name() string

	Generate(
		ctx context.Context,
		request Request,
	) (*Response, error)

	Stream(
		ctx context.Context,
		request Request,
		handler StreamHandler,
	) error
}
