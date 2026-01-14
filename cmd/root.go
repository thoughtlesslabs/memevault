package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var vaultFile string

var rootCmd = &cobra.Command{
	Use:   "envault",
	Short: "A portable, secure environment variable manager",
	Long: `Envault is a CLI tool to manage environment variables securely using Age encryption.
It supports sharing secrets via encrypted files and even hiding them inside memes.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&vaultFile, "vault", "secrets.jpg", "Path to the vault file (encrypted file or meme)")
}
