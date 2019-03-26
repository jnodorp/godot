package cmd

import (
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt [ciphertext]",
	Short: "Decrypt a secret",
	Long: `The encrypt command encrypts a decrypt which can then
be used in the dotfile repository.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Read password from standard input
		log.Info("enter password for decryption")

		response, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.WithError(err).Fatal("failed to read password")
		}

		plaintext, err := decrypt([]byte(args[0]), response)
		if err != nil {
			log.WithError(err).Fatal("failed to decrypt input")
		}

		log.Info(plaintext)
	},
}

func init() {
	RootCmd.AddCommand(decryptCmd)
}
