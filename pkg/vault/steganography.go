package vault

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"

	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
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

// MemeResponse struct for Imgflip API
type MemeResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Memes []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"memes"`
	} `json:"data"`
}

// FetchRandomMeme downloads a random meme from imgflip.
func FetchRandomMeme() ([]byte, string, error) {
	resp, err := http.Get("https://api.imgflip.com/get_memes")
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	var memeResp MemeResponse
	if err := json.NewDecoder(resp.Body).Decode(&memeResp); err != nil {
		return nil, "", err
	}

	if !memeResp.Success || len(memeResp.Data.Memes) == 0 {
		return nil, "", errors.New("failed to fetch memes from imgflip")
	}

	rand.Seed(time.Now().UnixNano())
	randomMeme := memeResp.Data.Memes[rand.Intn(len(memeResp.Data.Memes))]

	// Download the image
	imgResp, err := http.Get(randomMeme.URL)
	if err != nil {
		return nil, "", err
	}
	defer imgResp.Body.Close()

	imgData, err := io.ReadAll(imgResp.Body)
	if err != nil {
		return nil, "", err
	}

	return imgData, randomMeme.URL, nil
}
