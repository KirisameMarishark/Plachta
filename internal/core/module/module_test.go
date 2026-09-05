package module

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	root := t.TempDir()

	modulesDir := filepath.Join(root, "modules")

	if err := os.MkdirAll(filepath.Join(modulesDir, "dns"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(modulesDir, "reality"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(modulesDir, "firewall"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(modulesDir, "README"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := List(root)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"dns",
		"firewall",
		"reality",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected modules: got %v, want %v", got, want)
	}
}

func TestListMissingDirectory(t *testing.T) {
	root := t.TempDir()

	_, err := List(root)
	if err == nil {
		t.Fatal("expected error when modules directory does not exist")
	}
}

func TestInfoOf(t *testing.T) {
	root := t.TempDir()

	moduleDir := filepath.Join(root, "modules", "reality")

	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		t.Fatal(err)
	}

	files := []string{
		"README.md",
		"install.sh",
		"verify.sh",
	}

	for _, name := range files {
		if err := os.WriteFile(filepath.Join(moduleDir, name), []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	got, err := InfoOf(root, "reality")
	if err != nil {
		t.Fatal(err)
	}

	want := Info{
		Name:    "reality",
		Path:    moduleDir,
		Readme:  true,
		Install: true,
		Verify:  true,
		Config:  false,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected module info: got %+v, want %+v", got, want)
	}
}

func TestInfoOfMissingModule(t *testing.T) {
	root := t.TempDir()

	_, err := InfoOf(root, "missing")
	if err == nil {
		t.Fatal("expected error when module does not exist")
	}
}
