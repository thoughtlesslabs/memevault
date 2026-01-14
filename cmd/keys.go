package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage keys",
}

var keysShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show your public key",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		keyFile := filepath.Join(home, ".envault", "keys", "envault.key")
		
		content, err := os.ReadFile(keyFile)
		if err != nil {
			fmt.Printf("Error reading key file: %v\n", err)
			return
		}
		
		// Parse comment for public key
		// In a real app we'd parse the private key object to derive it, but we stored it in comment for convenience
		fmt.Println(string(content))
	},
}

func init() {
	rootCmd.AddCommand(keysCmd)
	keysCmd.AddCommand(keysShowCmd)
}
