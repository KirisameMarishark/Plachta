package module

import (
	"fmt"
	"os"
	"path/filepath"
)

type Info struct {
	Name    string
	Path    string
	Readme  bool
	Install bool
	Verify  bool
	Config  bool
}

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

func InfoOf(root, name string) (Info, error) {
	modulesDir := filepath.Join(root, "modules")
	moduleDir := filepath.Join(modulesDir, name)

	info, err := os.Stat(moduleDir)
	if err != nil {
		if os.IsNotExist(err) {
			return Info{}, fmt.Errorf("module %q not found", name)
		}

		return Info{}, fmt.Errorf("failed to inspect module %q: %w", name, err)
	}

	if !info.IsDir() {
		return Info{}, fmt.Errorf("module %q not found", name)
	}

	return Info{
		Name:    name,
		Path:    moduleDir,
		Readme:  fileExists(filepath.Join(moduleDir, "README.md")),
		Install: fileExists(filepath.Join(moduleDir, "install.sh")),
		Verify:  fileExists(filepath.Join(moduleDir, "verify.sh")),
		Config:  fileExists(filepath.Join(moduleDir, "config.sh")),
	}, nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
