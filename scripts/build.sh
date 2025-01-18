#!/bin/bash

# Build script for Server Status Reporter
set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Building Server Status Reporter...${NC}"
go build -o serverstatus

echo -e "${GREEN}Build complete!${NC}"
