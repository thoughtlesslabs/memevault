package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var grantCmd = &cobra.Command{
	Use:   "grant [NAME] [PUBLIC_KEY]",
	Short: "Grant access to another user (by public key)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		key := args[1]

		if keyFile == "" {
			home, _ := os.UserHomeDir()
			keyFile = filepath.Join(home, ".memevault", "keys", "memevault.key")
		}

		secrets, err := loadSecrets(vaultFile, keyFile)
		if err != nil {
			fmt.Printf("Error loading secrets: %v\n", err)
			return
		}

		recipients := getRecipients(secrets)

		// TODO: Check if name/key already exists?
		newRecipient := Recipient{Name: name, PublicKey: key}
		recipients = append(recipients, newRecipient)

		if err := saveSecrets(vaultFile, secrets, recipients); err != nil {
			fmt.Printf("Error granting access: %v\n", err)
			return
		}

		fmt.Printf("Granted access to %s (%s)\n", name, key)
	},
}

func init() {
	rootCmd.AddCommand(grantCmd)
	grantCmd.Flags().StringVar(&keyFile, "key", "", "Path to private key file")
}
