package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.2.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "agewrap",
		Short: "agewrap — minimal age wrapper (encrypt/decrypt/keygen)",
		Long:  "agewrap is a tiny wrapper around the system 'age' tools to make encrypt/decrypt/keygen easier for local workflows.",
	}

	rootCmd.Version = version

	// add subcommands
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(keygenCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}