package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var vaultFile string

const Version = "v1.2.0"

var rootCmd = &cobra.Command{
	Use:     "memevault",
	Version: Version,
	Short:   "A portable, secure environment variable manager (with memes)",
	Long: `Memevault is a CLI tool to manage environment variables securely using Age encryption and steganography.
It supports sharing secrets via encrypted files hidden inside memes.`,
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
