package subscription

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KirisameMarishark/Plachta/internal/core/subscription/providers"
)

const (
	subscriptionDir  = "/etc/plachta/subscription"
	subscriptionFile = "/etc/plachta/subscription/sub.txt"
)

type Provider interface {
	Name() string
	URI() (string, error)
}

type Generator struct {
	Providers []Provider
}

func New() Generator {
	return Generator{
		Providers: []Provider{
			providers.NewRealityProvider(),
		},
	}
}

func (g Generator) Generate() (string, error) {
	if err := os.MkdirAll(subscriptionDir, 0755); err != nil {
		return "", fmt.Errorf("create subscription directory: %w", err)
	}

	var uris []string

	for _, provider := range g.Providers {
		uri, err := provider.URI()
		if err != nil {
			return "", fmt.Errorf("provider %s: %w", provider.Name(), err)
		}

		uri = strings.TrimSpace(uri)

		if uri == "" {
			continue
		}

		uris = append(uris, uri)
	}

	content := strings.Join(uris, "\n")

	if content != "" {
		content += "\n"
	}

	if err := os.WriteFile(
		subscriptionFile,
		[]byte(content),
		0644,
	); err != nil {
		return "", fmt.Errorf("write subscription file: %w", err)
	}

	return filepath.Clean(subscriptionFile), nil
}
