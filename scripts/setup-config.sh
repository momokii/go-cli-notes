#!/usr/bin/env bash
# Configuration setup script for Knowledge Garden CLI

set -e

CONFIG_DIR="$HOME/.config/kg-cli"
CONFIG_FILE="$CONFIG_DIR/config.yaml"

# Cloud API URL (safe to hardcode - public endpoint)
CLOUD_API_URL="https://cli-notes-api.kelanach.xyz/"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "=== Knowledge Garden CLI - Configuration Setup ==="
echo ""

# Create config directory
mkdir -p "$CONFIG_DIR"

# Check if config already exists
if [ -f "$CONFIG_FILE" ]; then
    echo -e "${YELLOW}Configuration file already exists at:${NC}"
    echo "  $CONFIG_FILE"
    echo ""
    read -p "Do you want to overwrite it? [y/N]: " overwrite
    if [[ ! $overwrite =~ ^[Yy]$ ]]; then
        echo "Keeping existing configuration."
        exit 0
    fi
    echo ""
fi

echo "Choose your API setup:"
echo "  1) Local API (http://localhost:8080)"
echo "     - For self-hosted API running on your machine"
echo ""
echo "  2) Cloud API (https://cli-notes-api.kelanach.xyz/)"
echo "     - Use the hosted cloud service"
echo ""
echo "  3) Custom URL"
echo "     - Specify your own API endpoint"
echo ""
read -p "Enter choice [1-3, default: 1]: " choice

case $choice in
    1)
        API_URL="http://localhost:8080"
        echo -e "${GREEN}Selected: Local API${NC}"
        ;;
    2)
        API_URL="$CLOUD_API_URL"
        echo -e "${GREEN}Selected: Cloud API${NC}"
        ;;
    3)
        read -p "Enter your API URL: " API_URL
        # Remove trailing slash if present
        API_URL="${API_URL%/}"
        echo -e "${GREEN}Selected: Custom API ($API_URL)${NC}"
        ;;
    "")
        API_URL="http://localhost:8080"
        echo -e "${YELLOW}No valid choice selected, using default: Local API${NC}"
        ;;
esac

echo ""
echo -e "${BLUE}Creating configuration file...${NC}"

# Create config file
cat > "$CONFIG_FILE" <<EOF
api:
  base_url: "$API_URL"
  timeout: 30

editor:
  external_editor: "vi"

preferences:
  default_note_type: "note"
  auto_save_interval: 30
  theme: "dark"
EOF

echo ""
echo -e "${GREEN}Configuration saved successfully!${NC}"
echo ""
echo "Config file: $CONFIG_FILE"
echo "API URL: $API_URL"
echo ""
echo "You can always edit this file later to change your API URL."
echo ""
echo "To get started:"
echo "  1. Register: kg-cli register"
echo "  2. Login:   kg-cli login"
echo "  3. Status:   kg-cli status"
