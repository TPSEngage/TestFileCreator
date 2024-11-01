package main

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ulikunitz/xz"
)

const FFMPEG_VERSION = "4.4.1"

var ffmpegPath string

func EnsureFfmpegExists() error {
	if ffmpegPath != "" {
		return nil
	}

	executableName := "ffmpeg"
	if runtime.GOOS == "windows" {
		executableName += ".exe"
	}

	if _, err := os.Stat(executableName); err == nil {
		ffmpegPath = executableName
		return nil
	}

	err := downloadFFmpeg()
	if err != nil {
		return fmt.Errorf("failed to download FFmpeg: %v", err)
	}

	return nil
}

func GetFfmpegPath() string {
	return ffmpegPath
}

func downloadFFmpeg() error {
	fmt.Printf("Downloading FFmpeg %s...\n", FFMPEG_VERSION)

	var url string
	var archiveExt string

	switch runtime.GOOS {
	case "linux":
		url = "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"
		archiveExt = ".tar.xz"
	case "darwin":
		url = fmt.Sprintf("https://evermeet.cx/ffmpeg/ffmpeg-%s.zip", FFMPEG_VERSION)
		archiveExt = ".zip"
	case "windows":
		url = "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip"
		archiveExt = ".zip"
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tempFile, err := os.CreateTemp("", "ffmpeg-download-*"+archiveExt)
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return err
	}

	tempFile.Close()

	fmt.Println("Extracting FFmpeg...")

	switch runtime.GOOS {
	case "linux":
		return extractTarXz(tempFile.Name())
	case "darwin", "windows":
		return extractZip(tempFile.Name())
	}

	return nil
}

func extractTarXz(archivePath string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %v", err)
	}
	defer file.Close()

	xzReader, err := xz.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create xz reader: %v", err)
	}

	tarReader := tar.NewReader(xzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %v", err)
		}

		if strings.HasSuffix(header.Name, "ffmpeg") {
			outFile, err := os.OpenFile(filepath.Base(header.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return fmt.Errorf("failed to create output file: %v", err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("failed to extract file: %v", err)
			}

			ffmpegPath = outFile.Name()
			fmt.Printf("FFmpeg extracted to: %s\n", ffmpegPath)
			return nil
		}
	}

	return fmt.Errorf("ffmpeg binary not found in the archive")
}

func extractZip(archivePath string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, "ffmpeg") || strings.HasSuffix(file.Name, "ffmpeg.exe") {
			outFile, err := os.OpenFile(filepath.Base(file.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			rc, err := file.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			_, err = io.Copy(outFile, rc)
			if err != nil {
				return err
			}

			ffmpegPath = outFile.Name()
			fmt.Printf("FFmpeg extracted to: %s\n", ffmpegPath)
			return nil
		}
	}

	return fmt.Errorf("ffmpeg binary not found in the archive")
}

func CreateFFmpegCommand(inputs map[string]interface{}, filter string, output string, outputArgs map[string]interface{}) *exec.Cmd {
	args := []string{"-y"} // Force overwrite output files

	// Add input arguments first
	for k, v := range inputs {
		args = append(args, fmt.Sprintf("-%s", k), fmt.Sprint(v))
	}

	// Add filter if present
	if filter != "" {
		args = append(args, "-filter_complex", filter)
	}

	// Add output arguments
	for k, v := range outputArgs {
		args = append(args, fmt.Sprintf("-%s", k), fmt.Sprint(v))
	}

	// Add output file last
	args = append(args, output)

	return exec.Command(ffmpegPath, args...)
}
