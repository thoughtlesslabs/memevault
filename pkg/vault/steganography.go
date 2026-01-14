package vault

import (
	"bytes"
	"embed"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"io"
	"net/http"
	"os"
)

// Magic bytes to identify our payload at the end of an image
var MagicBytes = []byte("ENVAULT_MEME")

// Embed appends the encrypted payload to the image at imagePath.
// Format: [Original Image Bytes] [Payload] [Payload Length (8 bytes)] [Magic Bytes]
func Embed(imagePath string, payload []byte) error {
	f, err := os.OpenFile(imagePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write payload
	if _, err := f.Write(payload); err != nil {
		return err
	}

	// Write length of payload (int64, little endian)
	lengthBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(lengthBuf, uint64(len(payload)))
	if _, err := f.Write(lengthBuf); err != nil {
		return err
	}

	// Write Magic Bytes
	if _, err := f.Write(MagicBytes); err != nil {
		return err
	}

	return nil
}

// Extract reads the payload from the end of the image file.
func Extract(imagePath string) ([]byte, error) {
	f, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := stat.Size()
	minSize := int64(len(MagicBytes) + 8)
	if fileSize < minSize {
		return nil, errors.New("file too small to contain envault payload")
	}

	// Read Magic Bytes
	magicBuf := make([]byte, len(MagicBytes))
	if _, err := f.ReadAt(magicBuf, fileSize-int64(len(MagicBytes))); err != nil {
		return nil, err
	}

	if !bytes.Equal(magicBuf, MagicBytes) {
		// Fallback: maybe it's just a raw encrypted file?
		// For now, strict steganography check.
		return nil, errors.New("envault magic bytes not found in image")
	}

	// Read Length
	lengthBuf := make([]byte, 8)
	lengthPos := fileSize - int64(len(MagicBytes)) - 8
	if _, err := f.ReadAt(lengthBuf, lengthPos); err != nil {
		return nil, err
	}

	payloadLen := int64(binary.LittleEndian.Uint64(lengthBuf))
	if payloadLen <= 0 || payloadLen > fileSize-minSize {
		return nil, errors.New("invalid payload length detected")
	}

	// Read Payload
	payload := make([]byte, payloadLen)
	payloadPos := lengthPos - payloadLen
	if _, err := f.ReadAt(payload, payloadPos); err != nil {
		return nil, err
	}

	return payload, nil
}

// MemeResponse struct for meme-api.com
type MemeResponse struct {
	PostLink  string   `json:"postLink"`
	Subreddit string   `json:"subreddit"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	NSFW      bool     `json:"nsfw"`
	Spoiler   bool     `json:"spoiler"`
	Author    string   `json:"author"`
	Ups       int      `json:"ups"`
	Preview   []string `json:"preview"`
}

//go:embed memes/*.jpg
var embeddedMemes embed.FS

// FetchRandomMeme downloads a random meme from meme-api.com, or falls back to embedded memes.
func FetchRandomMeme() ([]byte, string, error) {
	// Try API first
	imgData, url, err := fetchFromAPI()
	if err == nil {
		return imgData, url, nil
	}

	// Fallback to embedded
	fmt.Printf("API failed (%v), using offline fallback...\n", err)
	return fetchEmbeddedMeme()
}

func fetchFromAPI() ([]byte, string, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("https://meme-api.com/gimme")
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	var meme MemeResponse
	if err := json.NewDecoder(resp.Body).Decode(&meme); err != nil {
		return nil, "", err
	}

	if meme.URL == "" {
		return nil, "", errors.New("failed to fetch meme from api")
	}

	// Download the image
	imgResp, err := client.Get(meme.URL)
	if err != nil {
		return nil, "", err
	}
	defer imgResp.Body.Close()

	imgData, err := io.ReadAll(imgResp.Body)
	if err != nil {
		return nil, "", err
	}

	return imgData, meme.URL, nil
}

func fetchEmbeddedMeme() ([]byte, string, error) {
	entries, err := embeddedMemes.ReadDir("memes")
	if err != nil {
		return nil, "", err
	}

	if len(entries) == 0 {
		return nil, "", errors.New("no embedded memes found")
	}

	rand.Seed(time.Now().UnixNano())
	chosen := entries[rand.Intn(len(entries))]

	data, err := embeddedMemes.ReadFile("memes/" + chosen.Name())
	if err != nil {
		return nil, "", err
	}

	return data, "offline-fallback://" + chosen.Name(), nil
}
