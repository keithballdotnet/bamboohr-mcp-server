#!/bin/bash

# Setup script for BambooHR MCP Server development environment
# This script installs the Git pre-commit hook to protect sensitive credentials

echo "Setting up BambooHR MCP Server development environment..."

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo "Error: This script must be run from the root of the git repository"
    exit 1
fi

# Install pre-commit hook
echo "Installing pre-commit hook..."
cp scripts/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

echo "✅ Pre-commit hook installed successfully!"
echo ""
echo "The pre-commit hook will automatically:"
echo "  - Replace BAMBOOHR_API_KEY with YOUR_API_KEY in commits"
echo "  - Replace BAMBOOHR_COMPANY with YOUR_COMPANY in commits"
echo "  - Keep your local .vscode/mcp.json file unchanged"
echo ""
echo "Setup complete! You can now safely commit changes without exposing sensitive credentials."
