package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [PATH]",
	Short: "Scan code for environment variable usage",
	Long: `Scans the specified directory (defaults to current) for common patterns 
indicating environment variable usage (e.g. os.Getenv("VAR"), process.env.VAR).`,
	Run: func(cmd *cobra.Command, args []string) {
		rootPath := "."
		if len(args) > 0 {
			rootPath = args[0]
		}

		fmt.Printf("Scanning %s for secrets...\n", rootPath)

		// Regex patterns for common languages
		// Note: These are heuristic simplifications.
		patterns := []*regexp.Regexp{
			// Go: os.Getenv("VAR")
			regexp.MustCompile(`os\.Getenv\(\s*["']([A-Z_][A-Z0-9_]*)["']\s*\)`),
			// JS/TS: process.env.VAR or process.env['VAR']
			regexp.MustCompile(`process\.env\.([A-Z_][A-Z0-9_]*)`),
			regexp.MustCompile(`process\.env\[\s*["']([A-Z_][A-Z0-9_]*)["']\s*\]`),
			// Python: os.environ.get('VAR') or os.environ['VAR']
			regexp.MustCompile(`os\.environ(?:\[\s*["']|.\s*get\(\s*["'])([A-Z_][A-Z0-9_]*)["']`),
		}

		foundVars := make(map[string]bool)

		err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip access errors
			}
			if info.IsDir() {
				// Basic ignore list
				name := info.Name()
				if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" {
					return filepath.SkipDir
				}
				return nil
			}

			// Skip likely binary/image files based on extension if needed, or just let regex fail fast
			// Simple check for text files?
			ext := filepath.Ext(path)
			if ext == ".exe" || ext == ".jpg" || ext == ".png" || ext == ".db" || ext == ".key" {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			for _, re := range patterns {
				matches := re.FindAllStringSubmatch(string(content), -1)
				for _, match := range matches {
					if len(match) > 1 {
						foundVars[match[1]] = true
					}
				}
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Error during scan: %v\n", err)
			return
		}

		if len(foundVars) == 0 {
			fmt.Println("No environment variables found in code.")
			return
		}

		fmt.Printf("Found %d potential variables:\n", len(foundVars))
		for v := range foundVars {
			fmt.Printf("- %s\n", v)
		}

		// Check against vault
		home, _ := os.UserHomeDir()
		if keyFile == "" {
			keyFile = filepath.Join(home, ".envault", "keys", "envault.key")
		}

		fmt.Println("\nChecking against vault...")
		secrets, err := loadSecrets(vaultFile, keyFile)
		missing := []string{}

		if err == nil {
			for v := range foundVars {
				if _, ok := secrets[v]; !ok {
					missing = append(missing, v)
				}
			}
		} else {
			fmt.Printf("Could not load vault to compare (%v). All found vars are potentially missing.\n", err)
			for v := range foundVars {
				missing = append(missing, v)
			}
		}

		if len(missing) > 0 {
			fmt.Printf("\n%d variables are MISSING from the vault:\n", len(missing))
			for _, v := range missing {
				fmt.Printf("[MISSING] %s\n", v)
			}
			fmt.Println("\nRun 'envault set <KEY> <VALUE>' to add them.")
		} else {
			fmt.Println("\nAll variables found in code are present in the vault. Good job!")
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVar(&keyFile, "key", "", "Path to private key file")
}
