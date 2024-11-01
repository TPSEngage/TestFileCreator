package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func escapeText(text string) string {
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, ":", "\\:")
	text = strings.ReplaceAll(text, "'", "\\'")
	return text
}

func createTestMedia(row map[string]string, fontPath string, outputDir string) error {
	resolution := row["Resolution"]
	contentType := row["Content Type"]
	duration := row["Duration"]

	isVideo := strings.ToLower(contentType) == "video"
	var ext string
	if isVideo {
		ext = "mp4"
	} else {
		ext = "jpg"
	}

	outputFileName := fmt.Sprintf("%s_%s_%s.%s", resolution, strings.ToLower(contentType), duration, ext)
	outputFile := filepath.Join(outputDir, outputFileName)

	var displayText string
	if isVideo {
		displayText = fmt.Sprintf("%s - Video - Duration: %ss", resolution, duration)
	} else {
		displayText = fmt.Sprintf("%s - Static - Duration: %ss", resolution, duration)
	}

	displayText = escapeText(displayText)

	drawtextFilter := fmt.Sprintf("drawtext=fontfile='%s':fontsize=32:fontcolor=white:x=(w-tw)/2:y=(h-th)/2:text='%s':box=1:boxcolor=black@0.5:boxborderw=5", fontPath, displayText)

	var cmd *exec.Cmd
	if isVideo {
		cmd = CreateFFmpegCommand(
			map[string]interface{}{"f": "lavfi", "i": "color=c=black:s=" + resolution + ":d=" + duration + ":r=25"},
			drawtextFilter,
			outputFile,
			map[string]interface{}{"c:v": "libx264", "t": duration, "pix_fmt": "yuv420p", "movflags": "+faststart"},
		)
	} else {
		cmd = CreateFFmpegCommand(
			map[string]interface{}{"f": "lavfi", "i": "color=c=black:s=" + resolution},
			drawtextFilter,
			outputFile,
			map[string]interface{}{"frames:v": "1"},
		)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error creating %s:\n", outputFileName)
		fmt.Println(stderr.String())
		return err
	}

	fmt.Printf("Created %s\n", outputFileName)
	return nil
}
