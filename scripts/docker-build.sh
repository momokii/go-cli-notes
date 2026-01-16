#!/bin/bash
# Docker build and push script with flags
# Usage: ./build.sh [-i image_name] [-v version] [-c]

# Exit immediately if a command exits with a non-zero status
set -e

# 1. Set Default Values
IMAGE_NAME="kelanach/go-cli-notes"
VERSION="latest"
BUILD_CLI=false # Default: CLI image is NOT built

# 2. Help function for documentation
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -i, --image <name>    Set image name (Default: $IMAGE_NAME)"
    echo "  -v, --version <tag>   Set image version (Default: $VERSION)"
    echo "  -c, --with-cli        Include CLI image build (Optional)"
    echo "  -h, --help            Show this help message"
    exit 1
}

# 3. Parse Arguments
while [ "$#" -gt 0 ]; do
    case $1 in
        -i|--image) 
            IMAGE_NAME="$2"
            shift 
            ;;
        -v|--version) 
            VERSION="$2"
            shift 
            ;;
        -c|--with-cli) 
            BUILD_CLI=true 
            ;;
        -h|--help) 
            usage 
            ;;
        *) 
            echo "Error: Unknown parameter passed: $1"
            usage 
            ;;
    esac
    shift
done

echo "=== Knowledge Garden - Docker Build ==="
echo "Image Base : $IMAGE_NAME"
echo "Version    : $VERSION"
echo "Build CLI? : $BUILD_CLI"
echo "---------------------------------------"

# --- Build API Image (Default / Mandatory) ---
echo "[1/2] Building API image..."
docker build -f ../Dockerfile.api -t "$IMAGE_NAME-api:$VERSION" ..

# --- Build CLI Image (Optional) ---
if [ "$BUILD_CLI" = true ]; then
    echo "[2/2] Building CLI image..."
    docker build -f ../Dockerfile.cli -t "$IMAGE_NAME-cli:$VERSION" ..
else
    echo "[2/2] Skipping CLI image build (use --with-cli to enable)..."
fi

# --- Tagging Logic ---
# Only run if version is NOT 'latest', to create the 'latest' alias
if [ "$VERSION" != "latest" ]; then
    echo ""
    echo "Tagging images as 'latest'..."
    
    # Tag API
    docker tag "$IMAGE_NAME-api:$VERSION" "$IMAGE_NAME-api:latest"
    
    # Tag CLI (Only if it was built)
    if [ "$BUILD_CLI" = true ]; then
        docker tag "$IMAGE_NAME-cli:$VERSION" "$IMAGE_NAME-cli:latest"
    fi
fi

# --- Summary ---
echo ""
echo "=== Build Complete ==="
echo "Images built successfully:"
echo "  ✅ $IMAGE_NAME-api:$VERSION"
if [ "$VERSION" != "latest" ]; then
    echo "  🏷️  $IMAGE_NAME-api:latest"
fi

if [ "$BUILD_CLI" = true ]; then
    echo "  ✅ $IMAGE_NAME-cli:$VERSION"
    if [ "$VERSION" != "latest" ]; then
        echo "  🏷️  $IMAGE_NAME-cli:latest"
    fi
fi