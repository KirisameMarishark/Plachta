package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/KirisameMarishark/Plachta/internal/core/cli"
	"github.com/KirisameMarishark/Plachta/internal/core/config"
	"github.com/KirisameMarishark/Plachta/internal/core/reality"
	"github.com/KirisameMarishark/Plachta/internal/core/subscription"
	"github.com/KirisameMarishark/Plachta/internal/core/system"
	"github.com/KirisameMarishark/Plachta/internal/core/version"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		cli.ShowHelp()
		return
	}

	switch args[0] {
	case "version", "-v", "--version":
		cli.ShowVersion(getVersion())

	case "help", "-h", "--help":
		cli.ShowHelp()

	case "system":
		showSystem()

	case "install":
		handleInstall(args)

	case "reality":
		handleReality(args)

	case "subscription":
		handleSubscription(args)

	case "config":
		if err := cli.HandleConfig(args, config.New()); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		if err := forwardToLegacyCLI(args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func getVersion() string {
	return strings.TrimSpace(version.Value)
}

func projectRoot() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}

	dir := filepath.Dir(exe)

	for {
		if _, err := os.Stat(filepath.Join(dir, "VERSION")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)

		if parent == dir {
			break
		}

		dir = parent
	}

	return "."
}

func showSystem() {
	info := system.Get()

	fmt.Printf("Operating System : %s\n", info.OS)
	fmt.Printf("Architecture     : %s\n", info.Architecture)
	fmt.Printf("Hostname         : %s\n", info.Hostname)

	if runtime.GOOS == "windows" {
		return
	}
}

func forwardToLegacyCLI(args []string) error {
	root := projectRoot()
	legacy := filepath.Join(root, "cmd", "plachta")

	if _, err := os.Stat(legacy); err != nil {
		return fmt.Errorf("legacy CLI not found: %w", err)
	}

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("bash", append([]string{legacy}, args...)...)
	} else {
		cmd = exec.Command(legacy, args...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func handleInstall(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  plachta install reality")
		return
	}

	switch args[1] {
	case "reality":
		if err := reality.Install(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println("Reality installation completed.")

	default:
		fmt.Println("Usage:")
		fmt.Println("  plachta install reality")
	}
}

func handleReality(args []string) {
	cfg := reality.New()

	if len(args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  plachta-go reality read")
		fmt.Println("  plachta-go reality verify")
		fmt.Println("  plachta-go reality uri")
		fmt.Println("  plachta-go reality generate")
		fmt.Println("  plachta-go reality install")
		return
	}

	switch args[1] {
	case "generate":
		uuid, err := reality.GenerateUUID()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		keyPair, err := reality.GenerateKeyPair()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		shortID, err := reality.GenerateShortID()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println("UUID       :", uuid)
		fmt.Println("PrivateKey :", keyPair.PrivateKey)
		fmt.Println("PublicKey  :", keyPair.PublicKey)
		fmt.Println("Hash32     :", keyPair.Hash32)
		fmt.Println("ShortID    :", shortID)

	case "install":
		if err := reality.Install(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println("Reality installation completed.")

	case "read":
		uuid, err := cfg.UUID()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		port, err := cfg.Port()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		serverName, err := cfg.ServerName()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		shortID, err := cfg.ShortID()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println("Reality Configuration")
		fmt.Println("------------------------------")
		fmt.Println("UUID       :", uuid)
		fmt.Println("Port       :", port)
		fmt.Println("ServerName :", serverName)
		fmt.Println("ShortID    :", shortID)

	case "verify":
		result := cfg.Verify()

		fmt.Println("Reality Verify")
		fmt.Println("------------------------------")

		for _, check := range result.Checks {
			if check.OK {
				fmt.Println("✔", check.Name)
			} else {
				fmt.Println("✘", check.Name)
			}
		}

		fmt.Println()
		fmt.Println("Result")
		fmt.Println("------------------------------")
		fmt.Println("Passed :", result.Passed())
		fmt.Println("Failed :", result.Failed())

		if result.URI != "" {
			fmt.Println()
			fmt.Println("Import URI")
			fmt.Println("------------------------------")
			fmt.Println(result.URI)
		}

		if !result.Valid() {
			os.Exit(1)
		}

	case "uri":
		uri, err := cfg.URI()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println(uri)

	default:
		fmt.Println("Usage:")
		fmt.Println("  plachta-go reality read")
		fmt.Println("  plachta-go reality verify")
		fmt.Println("  plachta-go reality uri")
	}
}

func handleSubscription(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  plachta-go subscription generate")
		return
	}

	switch args[1] {
	case "generate":
		path, err := subscription.New().Generate()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println("Subscription generated:")
		fmt.Println(path)

	default:
		fmt.Println("Usage:")
		fmt.Println("  plachta-go subscription generate")
	}
}
