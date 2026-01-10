#!/bin/bash
# Docker build and push script

set -e

IMAGE_NAME=${IMAGE_NAME:-"kelanach/knowledge-garden"}
VERSION=${VERSION:-"latest"}

echo "=== Knowledge Garden - Docker Build ==="
echo "Image: $IMAGE_NAME:$VERSION"
echo ""

# Build API image
echo "Building API image..."
docker build -f ../Dockerfile.api -t $IMAGE_NAME-api:$VERSION ..

# Build CLI image
echo "Building CLI image..."
docker build -f ../Dockerfile.cli -t $IMAGE_NAME-cli:$VERSION ..

# Tag as latest
docker tag $IMAGE_NAME-api:$VERSION $IMAGE_NAME-api:latest
docker tag $IMAGE_NAME-cli:$VERSION $IMAGE_NAME-cli:latest

echo ""
echo "=== Build Complete ==="
echo "Images built:"
echo "  - $IMAGE_NAME-api:$VERSION"
echo "  - $IMAGE_NAME-cli:$VERSION"
echo ""
echo "To push to registry:"
echo "  docker push $IMAGE_NAME-api:$VERSION"
echo "  docker push $IMAGE_NAME-cli:$VERSION"
