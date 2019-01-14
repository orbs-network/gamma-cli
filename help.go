package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

func commandShowHelp(requiredOptions []string) {
	fmt.Fprintf(os.Stderr, "Usage:\n\n")
	fmt.Fprintf(os.Stderr, "gamma-cli COMMAND [OPTIONS]\n\n")

	fmt.Fprintf(os.Stderr, "Commands:\n\n")
	sortedCommands := sortCommands()
	for _, name := range sortedCommands {
		cmd := commands[name]
		fmt.Fprintf(os.Stderr, "  %s %s %s\n", name, strings.Repeat(" ", 15-len(name)), cmd.desc)
		if cmd.args != "" {
			fmt.Fprintf(os.Stderr, "  %s  options: %s\n", strings.Repeat(" ", 15), cmd.args)
		}
		if cmd.example != "" {
			fmt.Fprintf(os.Stderr, "  %s  example: %s\n", strings.Repeat(" ", 15), cmd.example)
		}
		if cmd.example2 != "" {
			fmt.Fprintf(os.Stderr, "  %s           %s\n", strings.Repeat(" ", 15), cmd.example2)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}
	fmt.Fprintf(os.Stderr, "\n")

	fmt.Fprintf(os.Stderr, "Options:\n\n")
	showOptions()
	fmt.Fprintf(os.Stderr, "\n")

	fmt.Fprintf(os.Stderr, "Multiple environments (eg. local and testnet) can be defined in orbs-gamma-config.json configuration file.\n")
	fmt.Fprintf(os.Stderr, "See https://github.com/orbs-network/orbs-contract-sdk for more info.\n")
	fmt.Fprintf(os.Stderr, "\n")

	os.Exit(2)
}

func commandVersion(requiredOptions []string) {
	log("gamma-cli version v%s", GAMMA_CLI_VERSION)

	gammaVersion := verifyDockerInstalled()
	log("Gamma server version %s (docker)", gammaVersion)
}

func sortCommands() []string {
	res := make([]string, len(commands))
	for name, cmd := range commands {
		res[cmd.sort] = name
	}
	return res
}

// taken from package flag (func PrintDefaults)
func showOptions() {
	flag.VisitAll(func(f *flag.Flag) {
		// ignore list
		if strings.HasPrefix(f.Name, "arg") {
			return
		}

		s := fmt.Sprintf("  -%s", f.Name)
		name, usage := flag.UnquoteUsage(f)
		if len(name) > 0 {
			s += " " + name
		}
		if len(s) <= 4 {
			s += "\t"
		} else {
			s += "\n    \t"
		}
		s += strings.Replace(usage, "\n", "\n    \t", -1)

		if !isZeroValue(f, f.DefValue) {
			s += fmt.Sprintf(" (default %q)", f.DefValue)
		}
		fmt.Fprint(os.Stderr, s, "\n")
	})
}

// taken from package flag
func isZeroValue(f *flag.Flag, value string) bool {
	typ := reflect.TypeOf(f.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	if value == z.Interface().(flag.Value).String() {
		return true
	}
	switch value {
	case "false", "", "0":
		return true
	}
	return false
}
