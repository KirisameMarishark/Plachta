package cli

import "fmt"

func ShowHelp() {
	fmt.Println("Plachta")
	fmt.Println("Proxy Infrastructure Framework")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  plachta <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  install      Install Plachta modules")
	fmt.Println("  update       Update installed components")
	fmt.Println("  doctor       Diagnose system status")
	fmt.Println("  show         Show runtime information")
	fmt.Println("  verify       Verify configuration")
	fmt.Println("  backup       Backup configuration")
	fmt.Println("  subscription Manage subscriptions")
	fmt.Println("  config       Manage configuration")
	fmt.Println("  modules      List modules")
	fmt.Println("  module-info  Show module information")
	fmt.Println("  system       Show system information")
	fmt.Println("  version      Show version")
	fmt.Println("  help         Show this help")
}

func ShowVersion(version string) {
	fmt.Printf("Plachta %s\n", version)
}

func HandleConfig(args []string, cfg ConfigHandler) error {
	if len(args) < 2 || args[1] == "show" {
		data, err := cfg.Read()
		if err != nil {
			return err
		}

		fmt.Print(data)
		return nil
	}

	switch args[1] {
	case "get":
		if len(args) < 3 {
			return fmt.Errorf("usage: plachta config get <key>")
		}

		value, err := cfg.Get(args[2])
		if err != nil {
			return fmt.Errorf("configuration key not found: %s", args[2])
		}

		fmt.Println(value)
		return nil

	case "set":
		if len(args) < 4 {
			return fmt.Errorf("usage: plachta config set <key> <value>")
		}

		if err := cfg.Set(args[2], args[3]); err != nil {
			return err
		}

		fmt.Println("Configuration updated.")
		return nil

	default:
		return fmt.Errorf(
			"usage:\n  plachta config show\n  plachta config get <key>\n  plachta config set <key> <value>",
		)
	}
}

type ConfigHandler interface {
	Read() (string, error)
	Get(key string) (string, error)
	Set(key, value string) error
}
