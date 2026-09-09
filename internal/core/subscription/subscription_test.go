package subscription

import (
	"os"
	"path/filepath"
	"testing"
)

type testProvider struct {
	name string
	uri  string
	err  error
}

func (p testProvider) Name() string {
	return p.name
}

func (p testProvider) URI() (string, error) {
	return p.uri, p.err
}

func TestNew(t *testing.T) {
	s := New()

	if s == nil {
		t.Fatal("New() returned nil")
	}

	if s.OutputDir != "/etc/plachta/subscription" {
		t.Fatalf("OutputDir = %q, want %q", s.OutputDir, "/etc/plachta/subscription")
	}

	if len(s.Providers) != 1 {
		t.Fatalf("Providers length = %d, want 1", len(s.Providers))
	}
}

func TestGenerateWritesSubscriptionFile(t *testing.T) {
	dir := t.TempDir()

	s := New()
	s.OutputDir = dir
	s.Providers = []Provider{
		testProvider{
			name: "test",
			uri:  "vless://test@example.com:443",
		},
	}

	path, err := s.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	wantPath := filepath.Join(dir, "sub.txt")
	if path != wantPath {
		t.Fatalf("Generate() path = %q, want %q", path, wantPath)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	want := "vless://test@example.com:443\n"

	if string(data) != want {
		t.Fatalf("subscription file = %q, want %q", string(data), want)
	}
}

func TestGenerateSkipsProviderErrors(t *testing.T) {
	dir := t.TempDir()

	s := New()
	s.OutputDir = dir
	s.Providers = []Provider{
		testProvider{
			name: "broken",
			err:  os.ErrNotExist,
		},
		testProvider{
			name: "working",
			uri:  "vless://working@example.com:443",
		},
	}

	_, err := s.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "sub.txt"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	want := "vless://working@example.com:443\n"

	if string(data) != want {
		t.Fatalf("subscription file = %q, want %q", string(data), want)
	}
}
