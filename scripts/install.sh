#!/bin/bash

# Install script for Server Status Reporter
set -e

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
mkdir -p /opt/serverstatus
mkdir -p /var/log/serverstatus
mkdir -p /etc/serverstatus

# Copy files
echo -e "${YELLOW}Copying files...${NC}"
cp serverstatus /opt/serverstatus/
cp scripts/serverstatus.service /etc/systemd/system/
cp scripts/serverstatus.timer /etc/systemd/system/
cp scripts/logrotate.conf /etc/logrotate.d/serverstatus

# Set permissions
echo -e "${YELLOW}Setting permissions...${NC}"
chown -R root:root /opt/serverstatus
chmod 755 /opt/serverstatus
chmod 755 /opt/serverstatus/serverstatus
chmod 644 /etc/systemd/system/serverstatus.service
chmod 644 /etc/systemd/system/serverstatus.timer

# Configure environment
if [ ! -f "/etc/serverstatus/.env" ]; then
    echo -e "${YELLOW}Setting up environment configuration...${NC}"
    cp .env.example /etc/serverstatus/.env
    echo "Please edit /etc/serverstatus/.env with your email settings"
fi

# Reload systemd
echo -e "${YELLOW}Reloading systemd...${NC}"
systemctl daemon-reload
systemctl enable serverstatus.timer
systemctl start serverstatus.timer

echo -e "${GREEN}Installation complete!${NC}"
echo -e "${YELLOW}Please edit /etc/serverstatus/.env with your email settings${NC}"
echo -e "${YELLOW}The service will run daily at 8 AM EST${NC}"
