package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const FONT_URL = "https://github.com/google/fonts/raw/main/ufl/ubuntu/Ubuntu-Regular.ttf"
const FONT_NAME = "Ubuntu-Regular.ttf"

func downloadFont(fontPath string) error {
	fmt.Printf("Downloading %s...\n", FONT_NAME)
	resp, err := http.Get(FONT_URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(fontPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
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
		displayText = fmt.Sprintf("%s\\nVideo\\nDuration: %ss", resolution, duration)
	} else {
		displayText = fmt.Sprintf("%s\\nStatic\\nDuration: %ss", resolution, duration)
	}

	// Escape special characters in the display text
	displayText = escapeText(displayText)

	// Correct the drawtext filter by removing 'drawtext='
	drawtextFilter := fmt.Sprintf("fontfile='%s':fontsize=32:fontcolor=white:x=(w-tw)/2:y=(h-th)/2:text='%s':box=1:boxcolor=black@0.5:boxborderw=5", fontPath, displayText)

	var stderr bytes.Buffer

	if isVideo {
		err := ffmpeg.Input("color=c=black:s="+resolution+":d="+duration+":r=25", ffmpeg.KwArgs{"f": "lavfi"}).
			Filter("drawtext", ffmpeg.Args{drawtextFilter}).
			Output(outputFile, ffmpeg.KwArgs{"c:v": "libx264", "t": duration, "pix_fmt": "yuv420p", "movflags": "+faststart"}).
			OverWriteOutput().
			WithErrorOutput(&stderr).
			Run()
		if err != nil {
			fmt.Printf("Error creating %s:\n", outputFileName)
			fmt.Println(stderr.String())
			return err
		}
	} else {
		err := ffmpeg.Input("color=c=black:s="+resolution, ffmpeg.KwArgs{"f": "lavfi"}).
			Filter("drawtext", ffmpeg.Args{drawtextFilter}).
			Output(outputFile, ffmpeg.KwArgs{"frames:v": 1}).
			OverWriteOutput().
			WithErrorOutput(&stderr).
			Run()
		if err != nil {
			fmt.Printf("Error creating %s:\n", outputFileName)
			fmt.Println(stderr.String())
			return err
		}
	}

	fmt.Printf("Created %s\n", outputFileName)
	return nil
}

func main() {
	// Define command-line flags
	outputDirFlag := flag.String("o", "", "Output directory")
	flag.Parse()

	// Determine the output directory
	var outputDir string
	if *outputDirFlag != "" {
		outputDir = *outputDirFlag
	} else {
		// Use current working directory
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current working directory:", err)
			return
		}
		outputDir = cwd
	}

	// Ensure output directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		fmt.Printf("Output directory does not exist: %s\n", outputDir)
		return
	}

	// Handle positional argument for formats.csv
	var csvPath string
	args := flag.Args()
	if len(args) > 0 {
		csvPath = args[0]
		fmt.Println("Using formats.csv from:", csvPath)
	} else {
		// Get the directory of the executable
		execPath, err := os.Executable()
		if err != nil {
			fmt.Println("Error getting executable path:", err)
			return
		}
		execDir := filepath.Dir(execPath)
		csvPath = filepath.Join(execDir, "formats.csv")
		fmt.Println("Using formats.csv from executable directory:", csvPath)
	}

	// Ensure the CSV file exists
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("Error: formats.csv not found at", csvPath)
		return
	}

	// Ensure the font file exists in the output directory
	fontPath, err := ensureFontExists(outputDir)
	if err != nil {
		fmt.Println("Error ensuring font exists:", err)
		return
	}

	csvFile, err := os.Open(csvPath)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading CSV headers:", err)
		return
	}

	// Read all remaining records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV records:", err)
		return
	}

	for _, record := range records {
		row := make(map[string]string)
		for i, value := range record {
			row[headers[i]] = value
		}
		err := createTestMedia(row, fontPath, outputDir)
		if err != nil {
			fmt.Println("Error creating media:", err)
			return
		}
	}
}
