#!/bin/bash

set -e

echo "Installing Hostbook..."

if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go first."
    exit 1
fi

go install github.com/engnhn/hostbook@latest

if [ -f "$HOME/go/bin/hostbook" ]; then
    echo "Running hostbook setup..."
    "$HOME/go/bin/hostbook" setup
else
    echo "Error: Installation failed."
    exit 1
fi

echo "Hostbook installed successfully!"
