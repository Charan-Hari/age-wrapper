package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt [flags] <infile>",
	Short: "Encrypt a file for one or more age recipients",
	Args:  cobra.ExactArgs(1),
	RunE:  runEncrypt,
}

func init() {
	encryptCmd.Flags().StringArrayP("recipient", "r", []string{}, "recipient(s) (age public key string, github:username, or path to key file). Can be passed multiple times.")
	encryptCmd.Flags().StringP("out", "o", "", "output file (default: <infile>.age)")
	encryptCmd.Flags().BoolP("armor", "a", false, "ASCII-armored (pass --armor to 'age' if available)")
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	infile := args[0]
	recips, _ := cmd.Flags().GetStringArray("recipient")
	out, _ := cmd.Flags().GetString("out")
	armor, _ := cmd.Flags().GetBool("armor")

	if len(recips) == 0 {
		return errNoRecipients
	}

	// ensure 'age' exists
	if err := mustExecutable("age"); err != nil {
		return err
	}

	// resolve recipients (files -> contents); also handle github:username
	resolved := []string{}
	for _, r := range recips {
		r = strings.TrimSpace(r)
		if strings.HasPrefix(r, "github:") {
			user := strings.TrimPrefix(r, "github:")
			user = strings.TrimSpace(user)
			if user == "" {
				return fmt.Errorf("empty github username in recipient %q", r)
			}
			keys, err := fetchGitHubRecipients(user)
			if err != nil {
				return fmt.Errorf("fetching github keys for %s: %w", user, err)
			}
			resolved = append(resolved, keys...)
			continue
		}

		// otherwise resolve single recipient argument (may be a filepath)
		s, err := resolveRecipientArg(r)
		if err != nil {
			return fmt.Errorf("resolving recipient %q: %w", r, err)
		}
		resolved = append(resolved, s)
	}

	if len(resolved) == 0 {
		return errNoRecipients
	}

	// determine out file
	if out == "" {
		out = infile + ".age"
	}

	if err := ensureOutDir(out); err != nil {
		return err
	}

	// build args: age -o out [-a] -r RECIPIENT -r RECIPIENT infile
	argsAge := []string{"-o", out}
	if armor {
		argsAge = append(argsAge, "-a")
	}
	for _, r := range resolved {
		argsAge = append(argsAge, "-r", r)
	}
	argsAge = append(argsAge, infile)

	fmt.Println("agewrap: running: age", strings.Join(argsAge, " "))

	if err := runCmd("age", argsAge...); err != nil {
		return fmt.Errorf("encrypt failed: %w", err)
	}

	fmt.Println("Encrypted ->", out)
	return nil
}