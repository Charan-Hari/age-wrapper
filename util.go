package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// runCmd runs a command and captures combined output; returns error on failure.
func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command failed: %s %v\noutput:\n%s\n%w", name, args, out.String(), err)
	}
	return nil
}

// runCmdCapture runs command and returns combined output or error.
func runCmdCapture(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// fileExists returns true if path exists and is a regular file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// resolveRecipientArg: if arg is a file path, read its content and return the first non-empty line.
// otherwise returns arg unchanged.
func resolveRecipientArg(arg string) (string, error) {
	if fileExists(arg) {
		b, err := os.ReadFile(arg)
		if err != nil {
			return "", err
		}
		for _, line := range bytes.Split(bytes.TrimSpace(b), []byte{'\n'}) {
			line = bytes.TrimSpace(line)
			if len(line) > 0 {
				return string(line), nil
			}
		}
		return string(bytes.TrimSpace(b)), nil
	}
	// Not a file: treat as literal recipient (e.g., age1...)
	return arg, nil
}

// mustExecutable checks that exec name is available in PATH.
func mustExecutable(name string) error {
	_, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("required executable %q not found in PATH. Please install it and try again", name)
	}
	return nil
}

// ensureOutDir creates parent dir for output path if needed.
func ensureOutDir(out string) error {
	dir := filepath.Dir(out)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

// simple validation for recipient string
func isProbablyAgeRecipient(s string) bool {
	return len(s) > 4 && s[:4] == "age1"
}

// fetchGitHubRecipients fetches https://github.com/<username>.keys and returns the non-empty public key lines.
func fetchGitHubRecipients(user string) ([]string, error) {
	user = strings.TrimSpace(user)
	if user == "" {
		return nil, errors.New("empty github username")
	}
	url := "https://github.com/" + user + ".keys"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching github keys: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("github user %q not found (404)", user)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github keys fetch returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading github keys body: %w", err)
	}

	lines := strings.Split(string(body), "\n")
	out := []string{}
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		// keep lines that look like SSH public keys (ssh-*)
		if strings.HasPrefix(l, "ssh-") || strings.HasPrefix(l, "ecdsa-") || strings.HasPrefix(l, "ecdsa-sha2-") {
			out = append(out, l)
		} else {
			// also accept raw age recipients if present
			if isProbablyAgeRecipient(l) {
				out = append(out, l)
			}
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no usable public keys found for github:%s", user)
	}
	return out, nil
}

// convenience error
var errNoRecipients = errors.New("no recipients provided")