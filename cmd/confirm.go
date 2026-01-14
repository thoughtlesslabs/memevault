package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// askForConfirmation asks the user for confirmation. A user must type in "yes" or "y" and press enter.
// It returns true if the user confirmed, false otherwise.
func askForConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/N]: ", prompt)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" || response == "" {
			return false
		}
	}
}
