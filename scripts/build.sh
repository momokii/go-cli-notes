#!/usr/bin/env bash
# Build script for Knowledge Garden CLI
# Builds binaries for multiple platforms

set -e

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Change to the project root directory (parent of scripts/)
cd "$SCRIPT_DIR/.."

VERSION=${VERSION:-"dev"}
BUILD_DIR="build"
DIST_DIR="dist"

echo "=== Knowledge Garden CLI - Build Script ==="
echo "Version: $VERSION"
echo "Working directory: $(pwd)"
echo ""

# Create directories
rm -rf $BUILD_DIR $DIST_DIR
mkdir -p $BUILD_DIR
mkdir -p $DIST_DIR

# Build for multiple platforms
build_platform() {
    GOOS=$1
    GOARCH=$2
    OUTPUT="$BUILD_DIR/kg-cli-$GOOS-$GOARCH"

    # Add .exe for Windows
    if [ "$GOOS" = "windows" ]; then
        OUTPUT="$OUTPUT.exe"
    fi

    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="-X 'main.Version=$VERSION'" \
        -o $OUTPUT \
        ./cmd/cli

    # Create tarball/zip
    BASENAME=$(basename $OUTPUT)
    if [ "$GOOS" = "windows" ]; then
        (cd $BUILD_DIR && zip -q "../$DIST_DIR/${BASENAME}.zip" $BASENAME)
    else
        tar -czf "$DIST_DIR/${BASENAME}.tar.gz" -C $BUILD_DIR $BASENAME
    fi
}

# Build each platform
build_platform linux amd64
build_platform linux arm64
build_platform darwin amd64
build_platform darwin arm64
build_platform windows amd64

echo ""
echo "=== Build Complete ==="
echo "Version: $VERSION"
echo "Binaries available in: $DIST_DIR/"
ls -lh $DIST_DIR/
