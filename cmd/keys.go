package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		keyFile := filepath.Join(home, ".memevault", "keys", "memevault.key")

		content, err := os.ReadFile(keyFile)
		if err != nil {
			fmt.Printf("Error reading key file: %v\n", err)
			return
		}

		// Parse comment for public key
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "# Public Key: ") {
				fmt.Println(strings.TrimPrefix(line, "# Public Key: "))
				return
			}
		}
		fmt.Println("Error: Public key not found in key file")
	},
}

func init() {
	rootCmd.AddCommand(keysCmd)
	keysCmd.AddCommand(keysShowCmd)
}
