package model

import (
	"context"
	"testing"
)

type testProvider struct{}

func (testProvider) Name() string {
	return "test"
}

func (testProvider) Generate(
	ctx context.Context,
	request Request,
) (*Response, error) {
	return &Response{
		Content: "ok",
	}, nil
}

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	if err := registry.Register(
		testProvider{},
	); err != nil {
		t.Fatal(err)
	}

	provider, err := registry.Get("test")

	if err != nil {
		t.Fatal(err)
	}

	if provider.Name() != "test" {
		t.Fatalf(
			"unexpected provider: %s",
			provider.Name(),
		)
	}
}
