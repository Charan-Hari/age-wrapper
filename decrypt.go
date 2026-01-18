package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt [flags] <infile>",
	Short: "Decrypt a file using an identity (age secret key file)",
	Args:  cobra.ExactArgs(1),
	RunE:  runDecrypt,
}

func init() {
	decryptCmd.Flags().StringP("identity", "i", "", "path to identity (AGE secret key file) used to decrypt")
	decryptCmd.Flags().StringP("out", "o", "", "output file (default: input without .age)")
}

func runDecrypt(cmd *cobra.Command, args []string) error {
	infile := args[0]
	identity, _ := cmd.Flags().GetString("identity")
	out, _ := cmd.Flags().GetString("out")

	if identity == "" {
		return fmt.Errorf("identity (-i) is required")
	}

	// ensure 'age' exists
	if err := mustExecutable("age"); err != nil {
		return err
	}

	if !fileExists(identity) {
		return fmt.Errorf("identity file not found: %s", identity)
	}

	// default out file
	if out == "" {
		out = strings.TrimSuffix(infile, ".age")
		if out == infile {
			out = infile + ".dec"
		}
	}

	if err := ensureOutDir(out); err != nil {
		return err
	}

	// Build: age -d -i identity -o out infile
	argsAge := []string{"-d", "-i", identity, "-o", out, infile}

	fmt.Println("agewrap: running: age", strings.Join(argsAge, " "))

	if err := runCmd("age", argsAge...); err != nil {
		return fmt.Errorf("decrypt failed: %w", err)
	}

	fmt.Println("Decrypted ->", out)
	return nil
}