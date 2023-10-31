#!/bin/bash

# Function to check if the script is being run as root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        echo "Please run this script as root."
        exit 1
    fi
}

# Function to install packages with error handling
install_package() {
    package_name="$1"
    if ! dpkg -l | grep -q "^ii  $package_name "; then
        apt install -y "$package_name"
    else
        echo "$package_name is already installed."
    fi
}

# Function to download and install Go
install_go() {
    local go_version
    go_version=$(curl -sL https://golang.org/VERSION?m=text | head -1)
    local go_url="https://go.dev/dl/$go_version.linux-amd64.tar.gz"

    # Remove any existing Go installation
    rm -rf /usr/local/go

    # Download and extract the Go archive
    curl -sLo go.tar.gz "$go_url"
    tar -C /usr/local/ -xzf go.tar.gz
    rm go.tar.gz

    # Add Go binary path to system PATH
    echo "export PATH=\$PATH:/usr/local/go/bin" > /etc/profile.d/go.sh
    source /etc/profile.d/go.sh

    # Check installed Go version
    go version
}

# Main script

# Check if running as root
check_root

# Update and upgrade the system (Debian/Ubuntu)
apt update -q
apt upgrade -y

# Install dependencies for Go and Git
install_package "sudo"
install_package "curl"
install_package "build-essential"
install_package "ca-certificates"

# Install Go
install_go

# Install Git
install_package "git"

echo "Go and Git, along with their dependencies, have been successfully installed."
