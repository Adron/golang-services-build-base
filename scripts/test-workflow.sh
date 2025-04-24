#!/bin/bash

set -e  # Exit on error

echo "Testing GitHub Actions workflow locally..."

# Check if act is installed
if ! command -v act &> /dev/null; then
    echo "Error: act is not installed. Please install it first:"
    echo "  brew install act  # macOS"
    echo "  curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash  # Linux"
    echo "  choco install act-cli  # Windows"
    exit 1
fi

# Check if Docker is running
if ! docker info &> /dev/null; then
    echo "Error: Docker is not running. Please start Docker first."
    exit 1
fi

# Create a temporary directory for artifacts
mkdir -p .github/workflows/artifacts

# Run the workflow with M1/M2 architecture support
echo "Running workflow..."
act -l  # List available workflows

# Run specific jobs with architecture flag
echo "Running test job..."
act -j test --container-architecture linux/amd64 -v

echo "Running benchmark job..."
act -j benchmark --container-architecture linux/amd64 -v

echo "Running load-test job..."
act -j load-test --container-architecture linux/amd64 -v

# Clean up
rm -rf .github/workflows/artifacts

echo "Workflow testing completed" 