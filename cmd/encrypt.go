package cmd

import (
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt [plaintext]",
	Short: "Encrypt a secret",
	Long: `The encrypt command encrypts a secret which can then
be used in the dotfile repository.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Read password from standard input
		log.Info("enter password for encryption")

		response, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.WithError(err).Fatal("failed to read password")
		}

		ciphertext, err := encrypt([]byte(args[0]), response)
		if err != nil {
			log.WithError(err).Fatal("failed to encrypt input")
		}

		log.Info(ciphertext)
	},
}

func init() {
	RootCmd.AddCommand(encryptCmd)
}
