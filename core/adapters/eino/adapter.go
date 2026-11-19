package eino

import "github.com/fjk-dev/fjk-agent/core/model"

type Adapter struct {
	providers *model.Registry
}

func New(
	providers *model.Registry,
) *Adapter {
	return &Adapter{
		providers: providers,
	}
}
