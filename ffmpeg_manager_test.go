package main

import (
	"testing"
)

//func TestEnsureFfmpegExists(t *testing.T) {
//	// Clear ffmpegPath before the test
//	ffmpegPath = ""
//
//	err := EnsureFfmpegExists()
//	if err != nil {
//		t.Fatalf("EnsureFfmpegExists failed: %v", err)
//	}
//
//	if ffmpegPath == "" {
//		t.Fatal("ffmpegPath is empty after EnsureFfmpegExists")
//	}
//
//	// Check if the file actually exists
//	if _, err := os.Stat(ffmpegPath); os.IsNotExist(err) {
//		t.Fatalf("FFmpeg binary does not exist at path: %s", ffmpegPath)
//	}
//}

func TestGetFfmpegPath(t *testing.T) {
	// Set a dummy path
	ffmpegPath = "/dummy/path/ffmpeg"

	path := GetFfmpegPath()
	if path != ffmpegPath {
		t.Fatalf("GetFfmpegPath returned %s, expected %s", path, ffmpegPath)
	}
}

//func TestCreateFFmpegCommand(t *testing.T) {
//	// Set a dummy ffmpegPath for testing
//	ffmpegPath = "ffmpeg"
//
//	inputs := map[string]interface{}{
//		"f": "lavfi",
//		"i": "color=c=black:s=1280x720:d=5:r=30",
//	}
//	filter := "drawtext=fontfile='/path/to/font.ttf':fontsize=32:fontcolor=white:x=(w-tw)/2:y=(h-th)/2:text='Test'"
//	output := "output.mp4"
//	outputArgs := map[string]interface{}{
//		"c:v":      "libx264",
//		"t":        "5",
//		"pix_fmt":  "yuv420p",
//		"movflags": "+faststart",
//	}
//
//	cmd := CreateFFmpegCommand(inputs, filter, output, outputArgs)
//
//	expectedArgs := []string{
//		"-y",
//		"-f", "lavfi",
//		"-i", "color=c=black:s=1280x720:d=5:r=30",
//		"-filter_complex", "drawtext=fontfile='/path/to/font.ttf':fontsize=32:fontcolor=white:x=(w-tw)/2:y=(h-th)/2:text='Test'",
//		"-i", "output.mp4",
//		"-c:v", "libx264",
//		"-t", "5",
//		"-pix_fmt", "yuv420p",
//		"-movflags", "+faststart",
//	}
//
//	if cmd.Path != ffmpegPath {
//		t.Fatalf("Unexpected command path: got %s, want %s", cmd.Path, ffmpegPath)
//	}
//
//	for i, arg := range expectedArgs {
//		if i >= len(cmd.Args) || cmd.Args[i+1] != arg {
//			t.Fatalf("Unexpected argument at position %d: got %s, want %s", i, cmd.Args[i+1], arg)
//		}
//	}
//}

func TestDownloadFFmpeg(t *testing.T) {
	// This is a more complex test that actually downloads FFmpeg
	// It's commented out to avoid unnecessary downloads during routine testing
	// Uncomment and run manually when needed

	/*
		err := downloadFFmpeg()
		if err != nil {
			t.Fatalf("downloadFFmpeg failed: %v", err)
		}

		if ffmpegPath == "" {
			t.Fatal("ffmpegPath is empty after downloadFFmpeg")
		}

		// Check if the file actually exists
		if _, err := os.Stat(ffmpegPath); os.IsNotExist(err) {
			t.Fatalf("FFmpeg binary does not exist at path: %s", ffmpegPath)
		}

		// Clean up
		os.Remove(ffmpegPath)
	*/
}

//func TestExtractZip(t *testing.T) {
//	// Create a temporary zip file with a dummy FFmpeg executable
//	tempDir, err := os.MkdirTemp("", "ffmpeg-test")
//	if err != nil {
//		t.Fatalf("Failed to create temp directory: %v", err)
//	}
//	defer os.RemoveAll(tempDir)
//
//	zipPath := filepath.Join(tempDir, "test.zip")
//	dummyFFmpegPath := filepath.Join(tempDir, "ffmpeg")
//
//	// Create a dummy FFmpeg file
//	if err := os.WriteFile(dummyFFmpegPath, []byte("dummy ffmpeg"), 0755); err != nil {
//		t.Fatalf("Failed to create dummy FFmpeg: %v", err)
//	}
//
//	// Create a zip file containing the dummy FFmpeg
//	if err := createTestZip(zipPath, dummyFFmpegPath); err != nil {
//		t.Fatalf("Failed to create test zip: %v", err)
//	}
//
//	// Test extractZip
//	err = extractZip(zipPath)
//	if err != nil {
//		t.Fatalf("extractZip failed: %v", err)
//	}
//
//	// Check if FFmpeg was extracted
//	if _, err := os.Stat(ffmpegPath); os.IsNotExist(err) {
//		t.Fatalf("FFmpeg binary does not exist at path: %s", ffmpegPath)
//	}
//
//	// Clean up
//	os.Remove(ffmpegPath)
//}

func createTestZip(zipPath, filePath string) error {
	// Implementation of zip creation
	// This is a placeholder - you would need to implement this function
	// to create a test zip file containing a dummy FFmpeg executable
	return nil
}
