# Test Media Generator

A Go application that generates test media files (images and videos) based on a provided CSV configuration.

## Features

- Generates images and videos with specified resolutions and durations.
- Supports custom output directories.
- Cross-platform compatibility.

## Prerequisites

- [Go](https://golang.org/dl/) installed on your system.
- [FFmpeg](https://ffmpeg.org/download.html) installed and accessible in your PATH.

## Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/test-media-generator.git
cd test-media-generator
```

Build the application:

```bash
go build -o test_media_generator
```

Usage 

```bash
./test_media_generator -o /path/to/output formats.csv
```

