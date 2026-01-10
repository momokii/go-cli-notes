#!/usr/bin/env bash
# Installation script for Knowledge Garden CLI

set -e

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Change to the project root directory (parent of scripts/)
cd "$SCRIPT_DIR/.."

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=== Knowledge Garden CLI - Installation ==="
echo "Working directory: $(pwd)"
echo ""

# Detect OS
OS=$(uname -s)
ARCH=$(uname -m)

case $OS in
    Linux*)
        if [ "$ARCH" = "x86_64" ]; then
            BINARY="kg-cli-linux-amd64"
        elif [ "$ARCH" = "aarch64" ]; then
            BINARY="kg-cli-linux-arm64"
        else
            echo -e "${RED}Unsupported architecture: $ARCH${NC}"
            exit 1
        fi
        ;;
    Darwin*)
        if [ "$ARCH" = "x86_64" ]; then
            BINARY="kg-cli-darwin-amd64"
        elif [ "$ARCH" = "arm64" ]; then
            BINARY="kg-cli-darwin-arm64"
        else
            echo -e "${RED}Unsupported architecture: $ARCH${NC}"
            exit 1
        fi
        ;;
    *)
        echo -e "${RED}Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

INSTALL_DIR=${INSTALL_DIR:-"/usr/local/bin"}

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${YELLOW}Note: Installation requires root privileges${NC}"
    echo "Please run with sudo or specify INSTALL_DIR"
    echo ""
    echo "Example:"
    echo "  sudo ./scripts/install.sh"
    echo "  INSTALL_DIR=$HOME/.local/bin ./scripts/install.sh"
    exit 1
fi

# Check if binary exists in dist (from build.sh)
if [ ! -f "dist/$BINARY.tar.gz" ]; then
    echo -e "${RED}Error: Binary not found in dist/. Run ./scripts/build.sh first${NC}"
    exit 1
fi

# Extract and install
echo "Installing $BINARY to $INSTALL_DIR..."
tar -xzf "dist/$BINARY.tar.gz" -C /tmp
cp "/tmp/$BINARY" "$INSTALL_DIR/kg-cli"
chmod +x "$INSTALL_DIR/kg-cli"

# Create config directory
CONFIG_DIR="$HOME/.config/kg-cli"
mkdir -p "$CONFIG_DIR"

echo ""
echo -e "${GREEN}Installation successful!${NC}"
echo ""
echo "Binary installed to: $INSTALL_DIR/kg-cli"
echo "Config directory: $CONFIG_DIR"
echo ""
echo "To get started:"
echo "  1. Start the API server (or use Docker):"
echo "     docker compose up -d"
echo ""
echo "  2. Register an account:"
echo "     kg-cli register"
echo ""
echo "  3. Login:"
echo "     kg-cli login"
echo ""
echo "  4. Check status:"
echo "     kg-cli status"
