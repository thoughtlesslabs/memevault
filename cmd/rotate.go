package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jdiet/envault/pkg/vault"
	"github.com/spf13/cobra"
)

var keysRotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate your private key (re-encrypt vault and replace key)",
	Long: `Generates a new keypair, re-encrypts the vault to allow the new key 
and revoke the old one, and replaces your local key file (backing up the old one).`,
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		if keyFile == "" {
			keyFile = filepath.Join(home, ".envault", "keys", "envault.key")
		}

		// 1. Load current secrets with OLD key
		fmt.Println("Loading vault with current key...")
		secrets, err := loadSecrets(vaultFile, keyFile)
		if err != nil {
			fmt.Printf("Error loading secrets: %v\n", err)
			return
		}

		// 2. Load Old Public Key to identify it in the list
		oldKeyContent, err := os.ReadFile(keyFile)
		if err != nil {
			fmt.Printf("Error reading key file: %v\n", err)
			return
		}
		var oldPubKey string
		for _, line := range strings.Split(string(oldKeyContent), "\n") {
			if strings.HasPrefix(line, "# Public Key: ") {
				oldPubKey = strings.TrimSpace(strings.TrimPrefix(line, "# Public Key: "))
				break
			}
		}
		if oldPubKey == "" {
			fmt.Println("Could not determine old public key from file. Aborting.")
			return
		}

		// 3. Generate NEW Keypair
		fmt.Println("Generating new keypair...")
		newPriv, newPub, err := vault.GenerateKey()
		if err != nil {
			fmt.Printf("Error generating new key: %v\n", err)
			return
		}

		// 4. Update Recipients List
		currentRecipients := getRecipients(secrets)
		var newRecipients []string
		found := false
		for _, r := range currentRecipients {
			if r == oldPubKey {
				// Replace old with new
				newRecipients = append(newRecipients, newPub)
				found = true
			} else {
				// Keep others
				newRecipients = append(newRecipients, r)
			}
		}

		// If for some reason we weren't in the list (maybe single user implicit mode?), just add new.
		if !found {
			fmt.Println("Warning: Old key was not found in recipients list (maybe it was implicit?). Adding new key anyway.")
			newRecipients = append(newRecipients, newPub)
		}

		// 5. Save/Re-encrypt with NEW recipients
		// We need to update the internal recipients map too
		fmt.Println("Re-encrypting vault...")
		if err := saveSecrets(vaultFile, secrets, newRecipients); err != nil {
			fmt.Printf("Error saving vault: %v\n", err)
			return
		}

		// 6. Backup Old Key
		backupPath := keyFile + ".bak"
		fmt.Printf("Backing up old key to %s...\n", backupPath)
		if err := os.Rename(keyFile, backupPath); err != nil {
			fmt.Printf("Error backing up key: %v. \nCRITICAL: New key is NOT saved yet! New private key is:\n%s\nSave this manually!!\n", err, newPriv)
			return
		}

		// 7. Write New Key
		fmt.Println("Saving new key...")
		err = os.WriteFile(keyFile, []byte(newPriv+"\n# Public Key: "+newPub+"\n"), 0600)
		if err != nil {
			fmt.Printf("Error writing new key: %v.\nCRITICAL: Restore backup from %s immediately or save this private key:\n%s\n", err, backupPath, newPriv)
			// Try to revert backup?
			os.Rename(backupPath, keyFile)
			return
		}

		fmt.Println("Success! Key rotated.")
		fmt.Printf("New Public Key: %s\n", newPub)
	},
}

func init() {
	keysCmd.AddCommand(keysRotateCmd)
	keysRotateCmd.Flags().StringVar(&keyFile, "key", "", "Path to private key file")
}
