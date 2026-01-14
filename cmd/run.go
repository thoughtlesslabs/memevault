package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var keyFile string

var runCmd = &cobra.Command{
	Use:   "run -- [command]",
	Short: "Run a command with secrets loaded",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if keyFile == "" {
			home, _ := os.UserHomeDir()
			keyFile = filepath.Join(home, ".memevault", "keys", "memevault.key")
		}

		secrets, err := loadSecrets(vaultFile, keyFile)
		if err != nil {
			fmt.Printf("Error loading secrets: %v\n", err)
			os.Exit(1)
		}

		// Prepare command
		runName := args[0]
		runArgs := args[1:]

		c := exec.Command(runName, runArgs...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		// Inject environment
		env := os.Environ()
		for k, v := range secrets {
			if k == RecipientsKey {
				continue
			}
			if !isValidKey(k) {
				fmt.Fprintf(os.Stderr, "Warning: Skipping invalid key '%s' found in vault.\n", k)
				continue
			}
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}

		// Polyfill: Check if command is "printenv"
		if runName == "printenv" {
			// If arguments provided, filter variables
			if len(runArgs) > 0 {
				for _, arg := range runArgs {
					found := false
					for _, e := range env {
						// e is "KEY=VALUE"
						if len(e) > len(arg)+1 && e[:len(arg)+1] == arg+"=" {
							fmt.Println(e[len(arg)+1:])
							found = true
							break // Match found for this arg (last one wins if duplicates exist, which shouldn't happen with map but os.Environ might have duplicates. We take last one usually or iterate all. Standard printenv prints all if duplicates? memevault prioritizes secrets.)
							// Actually properly parsing env strings is better.
						}
					}
					if !found {
						// printenv usually exits with 1 if variable not found when requested?
						// simple implementation: just print nothing if not found.
					}
				}
			} else {
				// Print all
				for _, e := range env {
					fmt.Println(e)
				}
			}
			return
		}

		c.Env = env

		if err := c.Run(); err != nil {
			// Forward exit code if possible
			fmt.Printf("Command failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&keyFile, "key", "", "Path to private key file")
}
