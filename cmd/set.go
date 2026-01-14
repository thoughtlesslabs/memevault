package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [KEY] [VALUE]",
	Short: "Set a secret value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		val := args[1]

		if keyFile == "" {
			home, _ := os.UserHomeDir()
			keyFile = filepath.Join(home, ".memevault", "keys", "memevault.key")
		}

		secrets, err := loadSecrets(vaultFile, keyFile)
		if err != nil {
			fmt.Printf("Error loading secrets: %v\n", err)
			return
		}

		secrets[key] = val

		// Preserve existing recipients
		recipients := getRecipients(secrets)

		if err := saveSecrets(vaultFile, secrets, recipients); err != nil {
			fmt.Printf("Error saving secrets: %v\n", err)
			return
		}

		fmt.Printf("Set %s\n", key)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
