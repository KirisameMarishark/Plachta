package system

import (
	"runtime"
	"testing"
)

func TestHostname(t *testing.T) {
	name := hostname()

	if name == "" {
		t.Fatal("hostname returned empty string")
	}
}

func TestOSReleaseValue(t *testing.T) {
	value, err := osReleaseValue("ID")

	if err != nil {
		t.Skipf("cannot read /etc/os-release: %v", err)
	}

	if value == "" {
		t.Fatal("OS ID is empty")
	}
}

func TestCPUCount(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("nproc is only available on Linux")
	}

	count, err := cpuCount()

	if err != nil {
		t.Fatalf("cpuCount returned error: %v", err)
	}

	if count < 1 {
		t.Fatalf("unexpected CPU count: %d", count)
	}
}

func TestMemoryMB(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("free is only available on Linux")
	}

	memory, err := memoryMB()

	if err != nil {
		t.Fatalf("memoryMB returned error: %v", err)
	}

	if memory < 1 {
		t.Fatalf("unexpected memory size: %d MB", memory)
	}
}

func TestGet(t *testing.T) {
	info := Get()

	if info.Hostname == "" {
		t.Fatal("Get().Hostname is empty")
	}

	if info.Architecture == "" {
		t.Fatal("Get().Architecture is empty")
	}

	if info.OS == "" {
		t.Fatal("Get().OS is empty")
	}

	if runtime.GOOS != "windows" {
		if info.Version == "" {
			t.Fatal("Get().Version is empty")
		}

		if info.Kernel == "" {
			t.Fatal("Get().Kernel is empty")
		}

		if info.CPUCount < 1 {
			t.Fatalf("Get().CPUCount = %d, want >= 1", info.CPUCount)
		}

		if info.MemoryMB < 1 {
			t.Fatalf("Get().MemoryMB = %d, want >= 1", info.MemoryMB)
		}
	}
}
