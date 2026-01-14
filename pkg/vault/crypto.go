package vault

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"filippo.io/age"
)

// GenerateKey creates a new X25519 identity and its corresponding public key.
// Returns identity (private key) string and recipient (public key) string.
func GenerateKey() (string, string, error) {
	id, err := age.GenerateX25519Identity()
	if err != nil {
		return "", "", err
	}
	return id.String(), id.Recipient().String(), nil
}

// Encrypt encrypts the given data for the list of recipients.
func Encrypt(data []byte, recipients []string) ([]byte, error) {
	var parsedRecipients []age.Recipient
	for _, r := range recipients {
		pubKey, err := age.ParseX25519Recipient(r)
		if err != nil {
			return nil, fmt.Errorf("invalid recipient %q: %v", r, err)
		}
		parsedRecipients = append(parsedRecipients, pubKey)
	}

	out := &bytes.Buffer{}
	w, err := age.Encrypt(out, parsedRecipients...)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// Decrypt decrypts the given data using the provided identity (private key).
// We try to parse the identity string as an X25519 identity.
func Decrypt(data []byte, identityStr string) ([]byte, error) {
	identity, err := age.ParseX25519Identity(identityStr)
	if err != nil {
		// Try parsing as legacy if needed, or handle file paths?
		// For now assume direct string from key file
		return nil, fmt.Errorf("invalid identity: %v", err)
	}

	r, err := age.Decrypt(bytes.NewReader(data), identity)
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// LoadIdentityFromFile reads an identity key from a file.
// It skips comments and empty lines, returning the first valid identity found.
func LoadIdentityFromFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		return line, nil
	}
	return "", fmt.Errorf("no identity found in %s", path)
}
