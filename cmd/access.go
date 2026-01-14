package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var accessCmd = &cobra.Command{
	Use:   "access",
	Short: "Manage access to the vault",
	Long:  `List, add, or remove users who have access to the vault.`,
}

var accessListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users with access",
	Run: func(cmd *cobra.Command, args []string) {
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
		if len(recipients) == 0 {
			fmt.Println("No recipients found (implicit single-user mode or legacy vault).")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tPUBLIC KEY")
		for _, r := range recipients {
			fmt.Fprintf(w, "%s\t%s\n", r.Name, r.PublicKey)
		}
		w.Flush()
	},
}

var accessRemoveCmd = &cobra.Command{
	Use:   "remove [NAME|KEY]",
	Short: "Revoke access for a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		if keyFile == "" {
			home, _ := os.UserHomeDir()
			keyFile = filepath.Join(home, ".memevault", "keys", "memevault.key")
		}

		secrets, err := loadSecrets(vaultFile, keyFile)
		if err != nil {
			fmt.Printf("Error loading secrets: %v\n", err)
			return
		}

		// Safeguard: Do not remove self
		// Need current public key
		currentKey, err := os.ReadFile(keyFile)
		if err == nil {
			// Extract public key
			lines := strings.Split(string(currentKey), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "# Public Key: ") {
					myPub := strings.TrimSpace(strings.TrimPrefix(line, "# Public Key: "))
					if target == myPub {
						fmt.Println("Error: You cannot remove yourself!")
						return
					}
					// Also check if target name maps to my key
					recipients := getRecipients(secrets)
					for _, r := range recipients {
						if r.Name == target && r.PublicKey == myPub {
							fmt.Println("Error: You cannot remove yourself!")
							return
						}
					}
				}
			}
		}

		recipients := getRecipients(secrets)
		var newRecipients []Recipient
		removed := false

		for _, r := range recipients {
			if r.Name == target || r.PublicKey == target {
				removed = true
				continue
			}
			newRecipients = append(newRecipients, r)
		}

		if !removed {
			fmt.Printf("User '%s' not found.\n", target)
			return
		}

		if err := saveSecrets(vaultFile, secrets, newRecipients); err != nil {
			fmt.Printf("Error removing user: %v\n", err)
			return
		}

		fmt.Printf("Removed user '%s' and re-encrypted vault.\n", target)
	},
}

func init() {
	rootCmd.AddCommand(accessCmd)
	accessCmd.AddCommand(accessListCmd)
	accessCmd.AddCommand(accessRemoveCmd)
	// 'add' is handled by grant -> maybe we alias it later or move grant here?
	// For now keeping grant as top-level command but access add could be an alias.
}
