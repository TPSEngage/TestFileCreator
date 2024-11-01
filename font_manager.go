package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	FONT_URL  = "https://github.com/google/fonts/raw/main/ufl/ubuntu/Ubuntu-Regular.ttf"
	FONT_NAME = "Ubuntu-Regular.ttf"
)

func downloadFont(fontPath string) error {
	fmt.Printf("Downloading %s...\n", FONT_NAME)
	resp, err := http.Get(FONT_URL)
	if err != nil {
		return fmt.Errorf("failed to download font: %v", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(fontPath)
	if err != nil {
		return fmt.Errorf("failed to create font file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save font file: %v", err)
	}
	fmt.Println("Font downloaded successfully.")
	return nil
}

func ensureFontExists(fontDir string) (string, error) {
	fontPath := filepath.Join(fontDir, FONT_NAME)
	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		err := downloadFont(fontPath)
		if err != nil {
			return "", err
		}
	}
	return fontPath, nil
}
