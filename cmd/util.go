package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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
