#!/bin/bash

# Uninstall script for Server Status Reporter
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Stopping and disabling services...${NC}"
systemctl stop serverstatus.timer
systemctl disable serverstatus.timer
systemctl stop serverstatus.service
systemctl disable serverstatus.service

echo -e "${YELLOW}Removing files...${NC}"
rm -f /etc/systemd/system/serverstatus.service
rm -f /etc/systemd/system/serverstatus.timer
rm -f /etc/logrotate.d/serverstatus
rm -rf /opt/serverstatus
rm -rf /etc/serverstatus

echo -e "${YELLOW}Reloading systemd...${NC}"
systemctl daemon-reload

echo -e "${GREEN}Uninstallation complete!${NC}"
