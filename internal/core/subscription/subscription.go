package subscription

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KirisameMarishark/Plachta/internal/core/subscription/providers"
)

type Provider interface {
	Name() string
	URI() (string, error)
}

type Subscription struct {
	OutputDir string
	Providers []Provider
}

func New() *Subscription {
	return &Subscription{
		OutputDir: "/etc/plachta/subscription",
		Providers: ListProviders(),
	}
}

func ListProviders() []Provider {
	return []Provider{
		providers.NewRealityProvider(),
	}
}

func (s *Subscription) Generate() (string, error) {
	outputDir := s.OutputDir
	if outputDir == "" {
		outputDir = "/etc/plachta/subscription"
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("create subscription directory: %w", err)
	}

	outputFile := filepath.Join(outputDir, "sub.txt")

	var entries []string

	for _, provider := range s.Providers {
		uri, err := provider.URI()
		if err != nil {
			continue
		}

		uri = strings.TrimSpace(uri)
		if uri == "" {
			continue
		}

		entries = append(entries, uri)
	}

	content := ""
	if len(entries) > 0 {
		content = strings.Join(entries, "\n") + "\n"
	}

	if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("write subscription file: %w", err)
	}

	return outputFile, nil
}
