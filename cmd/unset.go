package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var forceUnset bool

var unsetCmd = &cobra.Command{
	Use:   "unset [KEY]",
	Short: "Remove a secret from the vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		if keyFile == "" {
			home, _ := os.UserHomeDir()
			keyFile = filepath.Join(home, ".memevault", "keys", "memevault.key")
		}

		secrets, err := loadSecrets(vaultFile, keyFile)
		if err != nil {
			fmt.Printf("Error loading secrets: %v\n", err)
			return
		}

		if _, ok := secrets[key]; !ok {
			fmt.Printf("Key '%s' not found in vault.\n", key)
			return
		}

		if !forceUnset {
			if !askForConfirmation(fmt.Sprintf("Are you sure you want to delete '%s'?", key)) {
				fmt.Println("Aborted.")
				return
			}
		}

		delete(secrets, key)

		// Preserve existing recipients
		recipients := getRecipients(secrets)

		if err := saveSecrets(vaultFile, secrets, recipients); err != nil {
			fmt.Printf("Error saving secrets: %v\n", err)
			return
		}

		fmt.Printf("Removed '%s'\n", key)
	},
}

func init() {
	rootCmd.AddCommand(unsetCmd)
	unsetCmd.Flags().BoolVarP(&forceUnset, "force", "f", false, "Skip confirmation prompt")
}
