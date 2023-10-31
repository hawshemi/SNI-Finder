#!/bin/bash


# Function to check if the script is being run as root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        echo 
        echo "Please run this script as root."
        echo 
        sleep 1
        exit 1
    fi
}


# Install Dependencies
install_package() {
    
    # Update and upgrade
    echo 
    echo "Updating the OS..."
    echo 
    sleep 0.5

    sudo apt update -q
    sudo apt upgrade -y

    echo 
    echo "OS Updated & Upgraded."
    echo 
    sleep 0.5

    echo 
    echo "Installing Dependencies..."
    echo 

    sudo apt install -y curl wget build-essential ca-certificates git

    echo 
    echo "Dependencies Installed."
    echo 
    sleep 0.5
}


# Function to download and install Go
install_go() {
    echo 
    echo "Installing Go..."
    echo 
    sleep 0.5
    local go_version
    go_version=$(curl -sL https://golang.org/VERSION?m=text | head -1)
    local go_url="https://go.dev/dl/$go_version.linux-amd64.tar.gz"
    sleep 0.5

    # Remove any existing Go installation
    sudo rm -rf /usr/local/go
    sleep 0.5

    # Download and extract the Go archive
    sudo curl -sLo go.tar.gz "$go_url"
    sleep 0.5
    
    sudo tar -C /usr/local/ -xzf go.tar.gz
    sleep 0.5

    sudo rm go.tar.gz

    # Add Go binary path to system PATH
    export PATH=$PATH:/usr/local/go/bin
    echo "export PATH=\$PATH:/usr/local/go/bin" > /etc/profile.d/go.sh
    sleep 0.5

    source /etc/profile.d/go.sh
    sleep 0.5

    # Check installed Go version
    go version
    sleep 0.5

    echo 
    echo "Go installed Succesfully."
    echo 
}


# Check if running as root
check_root


# Install dependencies for Go and Git
install_package


# Install Go
install_go



