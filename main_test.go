package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEscapeText(t *testing.T) {
	input := `Special characters: \ : '`
	expected := `Special characters: \\\\ \\\: \\\'`

	result := escapeText(input)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestEnsureFontExists(t *testing.T) {
	testDir := t.TempDir()
	fontPath, err := ensureFontExists(testDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		t.Errorf("Font file does not exist at path: %s", fontPath)
	}
}

func TestCreateTestMedia(t *testing.T) {
	testDir := t.TempDir()
	fontPath, err := ensureFontExists(testDir)
	if err != nil {
		t.Fatalf("Error ensuring font exists: %v", err)
	}

	row := map[string]string{
		"Resolution":   "640x480",
		"Content Type": "Static",
		"Duration":     "0",
	}

	err = createTestMedia(row, fontPath, testDir)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedFile := filepath.Join(testDir, "640x480_static_0.jpg")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist", expectedFile)
	}
}
