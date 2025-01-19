#!/bin/bash

# Install script for Server Status Reporter
set -e

INSTALL_DIR="/opt/serverstatus"
LOG_DIR="/var/log/serverstatus"
CONFIG_DIR="/etc/serverstatus"
SERVICE_DIR="/etc/systemd/system"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}Starting Server Status Reporter installation...${NC}"

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}Please run as root${NC}"
    exit 1
fi

# Create necessary directories
echo -e "${YELLOW}Creating directories...${NC}"
mkdir -p "$INSTALL_DIR"
mkdir -p "$LOG_DIR"
mkdir -p "$CONFIG_DIR"

# Copy files
echo -e "${YELLOW}Copying files...${NC}"
cp serverstatus "$INSTALL_DIR/"
cp scripts/serverstatus.service "$SERVICE_DIR/"
cp scripts/serverstatus.timer "$SERVICE_DIR/"
cp scripts/logrotate.conf /etc/logrotate.d/serverstatus

# Set permissions
echo -e "${YELLOW}Setting permissions...${NC}"
chown -R root:root "$INSTALL_DIR"
chmod 755 "$INSTALL_DIR"
chmod 755 "$INSTALL_DIR/serverstatus"
chmod 644 "$SERVICE_DIR/serverstatus.service"
chmod 644 "$SERVICE_DIR/serverstatus.timer"

# Configure environment
if [ ! -f "$CONFIG_DIR/.env" ]; then
    echo -e "${YELLOW}Setting up environment configuration...${NC}"
    cp .env.example "$CONFIG_DIR/.env"
    echo "Please edit $CONFIG_DIR/.env with your email settings"
fi

# Reload systemd
echo -e "${YELLOW}Reloading systemd...${NC}"
systemctl daemon-reload
systemctl enable serverstatus.timer
systemctl start serverstatus.timer

echo -e "${GREEN}Installation complete!${NC}"
echo -e "${YELLOW}Please edit $CONFIG_DIR/.env with your email settings${NC}"
echo -e "${YELLOW}The service will run daily at 8 AM EST${NC}"
