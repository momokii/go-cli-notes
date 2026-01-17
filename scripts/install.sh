#!/usr/bin/env bash
# Installation script for Knowledge Garden CLI
# Web-based installer that downloads pre-built binaries from GitHub Releases

set -e

# Configuration
REPO="momokii/go-cli-notes"
GITHUB_BASE_URL="https://github.com/$REPO"
CLOUD_API_URL="https://cli-notes-api.kelanach.xyz/"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "=== Knowledge Garden CLI - Installation ==="
echo ""

# Detect platform
OS=$(uname -s)
ARCH=$(uname -m)

case $OS in
    Linux*)
        PLATFORM="linux"
        ;;
    Darwin*)
        PLATFORM="darwin"
        ;;
    MINGW*|MSYS*|CYGWIN*)
        PLATFORM="windows"
        ;;
    *)
        echo -e "${RED}Unsupported OS: $OS${NC}"
        echo "This script supports Linux, macOS, and Windows (MSYS2/MinGW)"
        exit 1
        ;;
esac

case $ARCH in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    i386|i686)
        echo -e "${RED}Unsupported architecture: i386/i686 (32-bit)${NC}"
        echo "Please use a 64-bit system"
        exit 1
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${BLUE}Detected platform: $PLATFORM/$ARCH${NC}"
echo ""

# Binary name and download URL
BINARY_NAME="kg-cli-${PLATFORM}-${ARCH}"
if [ "$PLATFORM" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
    DOWNLOAD_URL="${GITHUB_BASE_URL}/releases/latest/download/${BINARY_NAME}.zip"
else
    DOWNLOAD_URL="${GITHUB_BASE_URL}/releases/latest/download/${BINARY_NAME}.tar.gz"
fi

# Install directory (prefer ~/.local/bin, fallback to /usr/local/bin)
INSTALL_DIR=""
NEEDS_SUDO=false

# Try ~/.local/bin first (user directory, no sudo required)
if mkdir -p "$HOME/.local/bin" 2>/dev/null; then
    if [ -w "$HOME/.local/bin" ]; then
        INSTALL_DIR="$HOME/.local/bin"
    fi
fi

# Fallback to /usr/local/bin (requires sudo)
if [ -z "$INSTALL_DIR" ]; then
    INSTALL_DIR="/usr/local/bin"
    NEEDS_SUDO=true
fi

echo -e "${BLUE}Installation directory: $INSTALL_DIR${NC}"
if [ "$NEEDS_SUDO" = true ]; then
    echo -e "${YELLOW}Note: sudo access required${NC}"
fi
echo ""

# Download and install
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

echo -e "${BLUE}Downloading $BINARY_NAME...${NC}"
if ! curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/archive.tar.gz"; then
    echo -e "${RED}Failed to download binary from GitHub Releases${NC}"
    echo ""
    echo "Possible reasons:"
    echo "  1. No releases have been published yet"
    echo "  2. Network connection issues"
    echo "  3. GitHub is experiencing issues"
    echo ""
    read -p "Would you like to install from source using 'go install' instead? [y/N]: " go_install < /dev/tty
    if [[ $go_install =~ ^[Yy]$ ]]; then
        echo ""
        echo -e "${BLUE}Installing from source using 'go install'...${NC}"
        echo "This may take a minute..."

        if [ "$NEEDS_SUDO" = true ]; then
            sudo go install "$GITHUB_BASE_URL/cmd/cli@latest" 2>/dev/null || \
            go install "$GITHUB_BASE_URL/cmd/cli@latest"
        else
            go install "$GITHUB_BASE_URL/cmd/cli@latest"
        fi

        if command -v kg-cli &> /dev/null; then
            echo ""
            echo -e "${GREEN}Installation from source successful!${NC}"
        else
            echo -e "${RED}Installation failed. Please make sure Go is installed.${NC}"
            exit 1
        fi
    else
        echo "Installation cancelled."
        exit 1
    fi
else
    # Extract archive
    echo -e "${BLUE}Extracting...${NC}"
    if [ "$PLATFORM" = "windows" ]; then
        unzip -q "$TMP_DIR/archive.tar.gz" -d "$TMP_DIR"
        BINARY_FILE="$TMP_DIR/${BINARY_NAME}"
    else
        tar -xzf "$TMP_DIR/archive.tar.gz" -C "$TMP_DIR"
        BINARY_FILE="$TMP_DIR/$BINARY_NAME"
    fi

    # Install binary
    echo -e "${BLUE}Installing to $INSTALL_DIR...${NC}"
    if [ "$NEEDS_SUDO" = true ]; then
        sudo cp "$BINARY_FILE" "$INSTALL_DIR/kg-cli"
    else
        cp "$BINARY_FILE" "$INSTALL_DIR/kg-cli"
    fi

    chmod +x "$INSTALL_DIR/kg-cli"
fi

# Verify installation
if command -v kg-cli &> /dev/null; then
    # Extract version from "kg-cli version v1.0.0" format
    VERSION=$(kg-cli --version 2>/dev/null | awk '{print $NF}' || echo "unknown")
    echo ""
    echo -e "${GREEN}Installation successful!${NC}"
    echo ""
    echo "Binary location: $INSTALL_DIR/kg-cli"
    echo "Version: $VERSION"
    echo ""
else
    echo -e "${RED}Installation verification failed${NC}"
    echo "The binary was installed but 'kg-cli' command is not found."
    echo "Make sure $INSTALL_DIR is in your PATH."
    exit 1
fi

# Setup configuration
echo -e "${BLUE}Setting up configuration...${NC}"
echo ""

# Check if config already exists
CONFIG_DIR="$HOME/.config/kg-cli"
CONFIG_FILE="$CONFIG_DIR/config.yaml"

if [ -f "$CONFIG_FILE" ]; then
    echo -e "${YELLOW}Configuration file already exists:${NC}"
    echo "  $CONFIG_FILE"
    echo ""
    read -p "Do you want to reconfigure? [y/N]: " reconfig < /dev/tty
    if [[ ! $reconfig =~ ^[Yy]$ ]]; then
        echo "Keeping existing configuration."
    else
        # Download and run setup config script
        SETUP_SCRIPT_URL="https://raw.githubusercontent.com/$REPO/main/scripts/setup-config.sh"
        if curl -fsSL "$SETUP_SCRIPT_URL" -o "$TMP_DIR/setup-config.sh"; then
            bash "$TMP_DIR/setup-config.sh"
        else
            echo -e "${YELLOW}Warning: Could not download setup script. You can configure manually later.${NC}"
        fi
    fi
else
    # Download and run setup config script
    SETUP_SCRIPT_URL="https://raw.githubusercontent.com/$REPO/main/scripts/setup-config.sh"
    if curl -fsSL "$SETUP_SCRIPT_URL" -o "$TMP_DIR/setup-config.sh"; then
        bash "$TMP_DIR/setup-config.sh"
    else
        echo -e "${YELLOW}Warning: Could not download setup script. You can configure manually later.${NC}"
    fi
fi

# Next steps
echo ""
echo -e "${GREEN}=== Installation Complete! ===${NC}"
echo ""
echo "To get started:"
echo "  1. Register an account:"
echo "     kg-cli register"
echo ""
echo "  2. Login:"
echo "     kg-cli login"
echo ""
echo "  3. Check status:"
echo "     kg-cli status"
echo ""
echo "For more information:"
echo "  https://github.com/$REPO"
echo ""
