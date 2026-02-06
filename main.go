package main

import (
	"fmt"
	"os"
)

// Version information (set by goreleaser at build time)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "sync":
		syncCommand(args)
	case "check":
		checkCommand(args)
	case "version", "--version", "-v":
		printVersion()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run 'gomodsync' without arguments for usage information")
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("gomodsync - Go module version comparison and sync tool")
	fmt.Println("\nUsage: gomodsync <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  sync       Synchronize dependency versions from reference to target")
	fmt.Println("  check      Check if target versions match reference")
	fmt.Println("  version    Show version information")
	fmt.Println("\nRun 'gomodsync <command> -h' for command-specific help")
}

func printVersion() {
	fmt.Printf("gomodsync version %s\n", version)
	fmt.Printf("  commit: %s\n", commit)
	fmt.Printf("  built:  %s\n", date)
}
