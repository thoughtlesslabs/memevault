package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thoughtlesslabs/memevault/pkg/vault"
)

type SecretsMap map[string]string

const RecipientsKey = "_memevault_recipients"

type Recipient struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

func getRecipients(secrets SecretsMap) []Recipient {
	val, ok := secrets[RecipientsKey]
	if !ok {
		return []Recipient{}
	}

	// Try parsing as new format
	var recipients []Recipient
	if err := json.Unmarshal([]byte(val), &recipients); err == nil {
		return recipients
	}

	// Fallback: Try parsing as old format ([]string)
	var keys []string
	if err := json.Unmarshal([]byte(val), &keys); err == nil {
		// Migrate
		migrated := make([]Recipient, len(keys))
		for i, k := range keys {
			migrated[i] = Recipient{Name: fmt.Sprintf("legacy-%s", k[:8]), PublicKey: k}
		}
		return migrated
	}

	return []Recipient{}
}

func setRecipients(secrets SecretsMap, recipients []Recipient) {
	data, _ := json.Marshal(recipients)
	secrets[RecipientsKey] = string(data)
}

func loadSecrets(vaultPath string, keyFile string) (SecretsMap, error) {
	// Read encrypted payload
	payload, err := vault.Extract(vaultPath)
	if err != nil {
		// Try reading as raw file if extraction fails (maybe not an image?)
		// For now simple fallback
		raw, err2 := os.ReadFile(vaultPath)
		if err2 != nil {
			return nil, fmt.Errorf("failed to read vault file: %v (also failed extract: %v)", err2, err)
		}
		payload = raw
	}

	// Decrypt
	identity, err := vault.LoadIdentityFromFile(keyFile)
	if err != nil {
		return nil, err
	}

	decrypted, err := vault.Decrypt(payload, identity)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %v", err)
	}

	var secrets SecretsMap
	if err := json.Unmarshal(decrypted, &secrets); err != nil {
		return nil, fmt.Errorf("invalid json payload: %v", err)
	}

	return secrets, nil
}

func saveSecrets(vaultPath string, secrets SecretsMap, recipients []Recipient) error {
	// Update internal metadata
	setRecipients(secrets, recipients)

	// Extract just keys for encryption
	var keys []string
	for _, r := range recipients {
		keys = append(keys, r.PublicKey)
	}

	data, err := json.Marshal(secrets)
	if err != nil {
		return err
	}

	encrypted, err := vault.Encrypt(data, keys)
	if err != nil {
		return err
	}

	// Check if vaultPath is an image
	ext := filepath.Ext(vaultPath)
	if ext == ".jpg" || ext == ".png" || ext == ".jpeg" {
		// It's an image, we need to re-embed.
		// NOTE: This implementation is destructive if we don't have the original image separately.
		// For a "vault" file that IS the image, we can't easily stripping the old payload without parsing.
		// Simpler approach for MVP: Expect the vaultPath to ALREADY be the image source.
		// But if we are UPDATING, we are just appending.
		// Wait, if I append to an image that already has a payload, I corrupt it or make it double payload.
		// FIX: We need a way to store the "clean" image or strip the payload.
		// IMPLEMENTATION SHORTCUT: We will not support "updating" the image in place perfectly without a clean source.
		// OR: We check for magic bytes, and if found, we truncate the file before appending new payload.

		f, err := os.OpenFile(vaultPath, os.O_RDWR, 0644)
		if err == nil {
			stat, _ := f.Stat()
			fileSize := stat.Size()

			// Check if it already has payload
			magicBuf := make([]byte, len(vault.MagicBytes))
			f.ReadAt(magicBuf, fileSize-int64(len(vault.MagicBytes)))

			if string(magicBuf) == string(vault.MagicBytes) {
				// It has a payload. Read length to find where to truncate.
				lengthBuf := make([]byte, 8)
				lengthPos := fileSize - int64(len(vault.MagicBytes)) - 8
				f.ReadAt(lengthBuf, lengthPos)
				// We don't even need the length, we just need to know it's there.
				// Actually, we need to find the START of the payload to truncate there.
				// But wait, the previous code doesn't store the start offset explicitly, it's calculated.
				// Let's rely on Extract to get the payload length, then calculate start.

				// Re-extract to find old payload size
				oldPayload, err := vault.Extract(vaultPath)
				if err == nil {
					// total overhead = len(oldPayload) + 8 (len) + len(magic)
					overhead := int64(len(oldPayload) + 8 + len(vault.MagicBytes))
					truncateAt := fileSize - overhead
					f.Truncate(truncateAt)
				}
			}
			f.Close()
		}

		return vault.Embed(vaultPath, encrypted)
	}

	// Normal file
	return os.WriteFile(vaultPath, encrypted, 0644)
}
