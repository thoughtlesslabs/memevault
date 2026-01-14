package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var grantCmd = &cobra.Command{
	Use:   "grant [PUBLIC_KEY]",
	Short: "Grant access to another user (by public key)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		newRecipient := args[0]

		if keyFile == "" {
			home, _ := os.UserHomeDir()
			keyFile = filepath.Join(home, ".memevault", "keys", "memevault.key")
		}

		secrets, err := loadSecrets(vaultFile, keyFile)
		if err != nil {
			fmt.Printf("Error loading secrets: %v\n", err)
			return
		}

		// Add new recipient to the list (logic handled in saveSecrets now)
		// We just pass it in. saveSecrets will merge it with existing ones.
		if err := saveSecrets(vaultFile, secrets, []string{newRecipient}); err != nil {
			fmt.Printf("Error granting access: %v\n", err)
			return
		}

		fmt.Printf("Granted access to %s\n", newRecipient)
	},
}

func init() {
	rootCmd.AddCommand(grantCmd)
	grantCmd.Flags().StringVar(&keyFile, "key", "", "Path to private key file")
}
