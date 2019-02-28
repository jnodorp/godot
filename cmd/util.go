package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

// askForConfirmation asks the user for confirmation. A user must type in "yes" or "no" and then press enter. It has
// fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as confirmations. If the input is not recognized, the
// default is used.
func userConfirm(s string, fallback bool) bool {
	reader := bufio.NewReader(os.Stdin)

	if fallback {
		fmt.Printf("%s [Y/n]: ", s)
	} else {
		fmt.Printf("%s [y/N]: ", s)
	}

	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	response = strings.ToLower(strings.TrimSpace(response))

	if response == "y" || response == "yes" {
		return true
	} else if response == "n" || response == "no" {
		return false
	} else {
		return fallback
	}
}

// expand tilde in the configured location.
func expand(s string) string {
	dir, err := homedir.Expand(s)
	if err != nil {
		log.Fatal("failed to expanding users home directory", err)
	}

	return dir
}

func homeDir() string {
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal("failed to determine users home directory", err)
	}

	return dir
}

// copy a file preserving its mode.
func copyFile(src, dst string) error {
	// Determine file mode.
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		in.Close()
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	in.Close()
	if err != nil {
		return err
	}

	return nil
}
