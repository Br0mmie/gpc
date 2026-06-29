#!/usr/bin/env bash

echo "This will erase any existing installation of gpc. Is this understood? "
read -p "[y/n]: " answer

if [[ "$answer" =~ ^[Yy]$ ]]; then
	echo "Removing old versions and installing new one..."
else
	echo "Installation cancelled"
	exit 1
fi

set -e

echo "Removing previous installs if any is found..."
sudo rm -f /usr/local/bin/gpc

echo "Creating temporary directory..."
TMP_DIR=$(mktemp -d)
echo "Location of temp directory: $TMP_DIR"

echo "Cloning git repo..."
sudo git clone --quiet https://github.com/Br0mmie/gpc.git "$TMP_DIR"
cd "$TMP_DIR"

echo "Building go files to /usr/local/bin/gpc..."
sudo go build -o /usr/local/bin/gpc -buildvcs=false > /dev/null 2>&1

echo "Chmodding executable..."
sudo chmod +x /usr/local/bin/gpc

if [ -f "/usr/local/bin/gpc" ] && [ -x "/usr/local/bin/gpc" ]; then
	echo "Removing temp directory $TMP_DIR"
	sudo rm -rf "$TMP_DIR"
    	echo "Foge has been installed!"
else
	echo "Error: Failed to install foge to /usr/local/bin/gpc"
    	echo "Removing temp directory $TMP_DIR"
	sudo rm -rf "$TMP_DIR"
    	exit 1
fi
