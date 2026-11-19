package model

import "context"

type Provider interface {
	Name() string

	Generate(
		ctx context.Context,
		request Request,
	) (*Response, error)
}
