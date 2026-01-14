package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
			keyFile = filepath.Join(home, ".envault", "keys", "envault.key")
		}

		secrets, err := loadSecrets(vaultFile, keyFile)
		if err != nil {
			fmt.Printf("Error loading secrets: %v\n", err)
			return
		}

		secrets[key] = val

		// For MVP, we only re-encrypt for the current user (ourselves).
		// In a real multi-user system, we need to know ALL recipients.
		// We'd store the recipient list in cleartext metadata in the vault file or assume config.
		// For this "Team Handover" demo, we will rely on re-encrypting for ourselves.
		// TODO: Store recipients in the vault header or dedicated file.
		
		// Hack for demo: Read public key from keyfile
		keyContent, _ := os.ReadFile(keyFile)
		lines := strings.Split(string(keyContent), "\n")
		var recipient string
		for _, line := range lines {
			if strings.HasPrefix(line, "# Public Key: ") {
				recipient = strings.TrimPrefix(line, "# Public Key: ")
				break
			}
		}

		if err := saveSecrets(vaultFile, secrets, []string{recipient}); err != nil {
			fmt.Printf("Error saving secrets: %v\n", err)
			return
		}

		fmt.Printf("Set %s\n", key)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
