package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thoughtlesslabs/memevault/pkg/vault"
)

var useMeme bool
var sourceImage string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new vault",
	Long:  `Create a new keypair and an empty vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Setup Keys Directory
		home, _ := os.UserHomeDir()
		keyDir := filepath.Join(home, ".memevault", "keys")
		if err := os.MkdirAll(keyDir, 0700); err != nil {
			fmt.Printf("Error creating key dir: %v\n", err)
			return
		}

		// 2. Load or Generate Keypair
		var pub string
		keyPath := filepath.Join(keyDir, "memevault.key")

		if _, err := os.Stat(keyPath); err == nil {
			fmt.Printf("Key already exists at %s. Using existing key.\n", keyPath)
			// Load public key from file (stored as comment)
			content, err := os.ReadFile(keyPath)
			if err != nil {
				fmt.Printf("Error reading existing key: %v\n", err)
				return
			}
			// Parse comment for public key
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "# Public Key: ") {
					pub = strings.TrimSpace(strings.TrimPrefix(line, "# Public Key: "))
					break
				}
			}
			if pub == "" {
				fmt.Println("Could not find public key in existing key file. Please backup and remove it to re-init.")
				return
			}
		} else {
			priv, newPub, err := vault.GenerateKey()
			if err != nil {
				fmt.Printf("Error generating key: %v\n", err)
				return
			}
			pub = newPub

			err = os.WriteFile(keyPath, []byte(priv+"\n# Public Key: "+pub+"\n"), 0600)
			if err != nil {
				fmt.Printf("Error checking writing key: %v\n", err)
				return
			}
			fmt.Printf("Created identity: %s\n", keyPath)
			fmt.Printf("Public Key: %s\n", pub)
		}

		// 3. Create Vault
		initialSecrets := SecretsMap{"Example": "Welcome to Envault"}

		// If meme requested or default to meme if no generic file
		finalVaultPath := vaultFile

		if useMeme || sourceImage == "" {
			// Fetch meme
			fmt.Println("Fetching a fresh meme from the internet...")
			imgData, url, err := vault.FetchRandomMeme()
			if err != nil {
				fmt.Printf("Failed to fetch meme: %v\n", err)
				return
			}
			fmt.Printf("Got meme: %s\n", url)

			// Save raw image first
			if filepath.Ext(finalVaultPath) == "" {
				finalVaultPath = "secrets.jpg"
			}
			if err := os.WriteFile(finalVaultPath, imgData, 0644); err != nil {
				fmt.Printf("Failed to write image: %v\n", err)
				return
			}
		} else {
			// Copy source image
			data, err := os.ReadFile(sourceImage)
			if err != nil {
				fmt.Printf("Failed to read source image: %v\n", err)
				return
			}
			finalVaultPath = filepath.Base(sourceImage) // simplistic
			os.WriteFile(finalVaultPath, data, 0644)
		}

		// Save secrets
		if err := saveSecrets(finalVaultPath, initialSecrets, []string{pub}); err != nil {
			fmt.Printf("Error creating vault: %v\n", err)
			return
		}

		fmt.Printf("Initialized vault at %s\n", finalVaultPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&useMeme, "meme", false, "Use a random meme as the vault container")
	initCmd.Flags().StringVar(&sourceImage, "image", "", "Use a specific image file as the vault container")
}
