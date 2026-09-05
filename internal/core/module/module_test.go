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
