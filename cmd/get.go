package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [KEY]",
	Short: "Get a secret value or list all secrets",
	Long:  `Retrieve a specific secret by key, or list all secrets if no key is provided.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if keyFile == "" {
			home, _ := os.UserHomeDir()
			keyFile = filepath.Join(home, ".memevault", "keys", "memevault.key")
		}

		// loadSecrets inherently checks access because it attempts to decrypt
		// with the user's private key. If they don't have access, this returns error.
		secrets, err := loadSecrets(vaultFile, keyFile)
		if err != nil {
			fmt.Printf("Error loading secrets: %v\n", err)
			return
		}

		if len(args) > 0 {
			// Get specific key
			key := args[0]
			if val, ok := secrets[key]; ok {
				fmt.Println(val)
			} else {
				fmt.Printf("Secret '%s' not found.\n", key)
				os.Exit(1)
			}
		} else {
			// List all keys
			// Sort keys for consistent output
			var keys []string
			for k := range secrets {
				if k == RecipientsKey {
					continue
				}
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Printf("%s=%s\n", k, secrets[k])
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
