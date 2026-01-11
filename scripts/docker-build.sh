#!/bin/bash
# Docker build and push script with flags

set -e

# 1. Set Default Values
IMAGE_NAME="kelanach/go-cli-notes"
VERSION="latest"

# 2. Fungsi Help untuk dokumentasi penggunaan
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -i, --image <name>    Set image name (Default: kelanach/go-cli-notes)"
    echo "  -v, --version <tag>   Set image version (Default: latest)"
    echo "  -h, --help            Show this help message"
    exit 1
}

# 3. Parsing Arguments (Looping cek setiap parameter)
while [ "$#" -gt 0 ]; do
    case $1 in
        -i|--image) IMAGE_NAME="$2"; shift ;;
        -v|--version) VERSION="$2"; shift ;;
        -h|--help) usage ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

echo "=== Knowledge Garden - Docker Build ==="
echo "Image: $IMAGE_NAME"
echo "Tag  : $VERSION"
echo ""

# Build API image
echo "Building API image..."
docker build -f ../Dockerfile.api -t "$IMAGE_NAME-api:$VERSION" ..

# Build CLI image
echo "Building CLI image..."
docker build -f ../Dockerfile.cli -t "$IMAGE_NAME-cli:$VERSION" ..

# Tag as latest (Opsional: biasanya hanya dilakukan jika version bukan latest)
if [ "$VERSION" != "latest" ]; then
    echo "Tagging as latest..."
    docker tag "$IMAGE_NAME-api:$VERSION" "$IMAGE_NAME-api:latest"
    docker tag "$IMAGE_NAME-cli:$VERSION" "$IMAGE_NAME-cli:latest"
fi

echo ""
echo "=== Build Complete ==="
echo "Images built:"
echo "  - $IMAGE_NAME-api:$VERSION"
echo "  - $IMAGE_NAME-cli:$VERSION"