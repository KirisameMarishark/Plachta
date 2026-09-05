package module

import (
	"fmt"
	"os"
	"path/filepath"
)

func List(root string) ([]string, error) {
	modulesDir := filepath.Join(root, "modules")

	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read modules directory: %w", err)
	}

	modules := make([]string, 0)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		modules = append(modules, entry.Name())
	}

	return modules, nil
}
