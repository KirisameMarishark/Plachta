package providers

import (
	"fmt"

	"github.com/KirisameMarishark/Plachta/internal/core/reality"
)

type RealityProvider struct {
	Config reality.Config
}

func NewRealityProvider() RealityProvider {
	return RealityProvider{
		Config: reality.New(),
	}
}

func (p RealityProvider) Name() string {
	return "reality"
}

func (p RealityProvider) URI() (string, error) {
	uri, err := p.Config.URI()
	if err != nil {
		return "", fmt.Errorf("generate reality URI: %w", err)
	}

	return uri, nil
}
