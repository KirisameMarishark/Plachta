package system

import (
	"os"
	"runtime"
)

type Info struct {
	OS           string
	Architecture string
	Hostname     string
}

func Get() Info {
	return Info{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		Hostname:     hostname(),
	}
}

func hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return "unknown"
	}

	return name
}
