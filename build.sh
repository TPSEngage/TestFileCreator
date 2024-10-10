#!/bin/bash

# Set the output directory
OUTPUT_DIR="dist"
mkdir -p $OUTPUT_DIR

# Disable CGO for cross-compilation
export CGO_ENABLED=0

# List of platforms to build for
PLATFORMS=(
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
)

# Build the project for each platform
for PLATFORM in "${PLATFORMS[@]}"; do
    OS=$(echo $PLATFORM | cut -d'/' -f1)
    ARCH=$(echo $PLATFORM | cut -d'/' -f2)
    OUTPUT_NAME="test_media_generator"

    if [ "$OS" = "windows" ]; then
        OUTPUT_NAME+=".exe"
    fi

    echo "Building for $OS/$ARCH..."

    # Set environment variables for cross-compilation
    env GOOS=$OS GOARCH=$ARCH go build -o $OUTPUT_DIR/$OS-$ARCH/$OUTPUT_NAME

    if [ $? -ne 0 ]; then
        echo "An error occurred while building for $OS/$ARCH."
        exit 1
    fi
done

echo "Builds completed successfully."
