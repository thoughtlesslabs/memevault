package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jdiet/envault/pkg/vault"
	"github.com/spf13/cobra"
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
		keyDir := filepath.Join(home, ".envault", "keys")
		if err := os.MkdirAll(keyDir, 0700); err != nil {
			fmt.Printf("Error creating key dir: %v\n", err)
			return
		}

		// 2. Generate Keypair
		priv, pub, err := vault.GenerateKey()
		if err != nil {
			fmt.Printf("Error generating key: %v\n", err)
			return
		}

		keyPath := filepath.Join(keyDir, "envault.key")
		// Don't overwrite existing
		if _, err := os.Stat(keyPath); err == nil {
			fmt.Printf("Key already exists at %s. backups recommended before re-init.\n", keyPath)
			// in a real app check force flag
		} else {
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

		if useMeme {
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
		} else if sourceImage != "" {
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
