package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var keygenCmd = &cobra.Command{
	Use:   "keygen [flags] <out-file>",
	Short: "Generate a new age keypair (wraps age-keygen)",
	Args:  cobra.ExactArgs(1),
	RunE:  runKeygen,
}

func init() {
	keygenCmd.Flags().BoolP("show-public", "p", false, "print the public key after generation")
}

func runKeygen(cmd *cobra.Command, args []string) error {
	outfile := args[0]
	showpub, _ := cmd.Flags().GetBool("show-public")

	// ensure age-keygen exists
	if err := mustExecutable("age-keygen"); err != nil {
		return err
	}

	// run: age-keygen -o outfile
	fmt.Println("agewrap: running: age-keygen -o", outfile)
	if err := runCmd("age-keygen", "-o", outfile); err != nil {
		return fmt.Errorf("keygen failed: %w", err)
	}
	fmt.Println("Wrote identity to", outfile)

	if showpub {
		// read file and print line with 'public key'
		f, err := os.Open(outfile)
		if err != nil {
			return err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "# public key:") {
				fmt.Println(line)
			}
		}
		// If nothing printed, show full file (last resort)
	}

	return nil
}