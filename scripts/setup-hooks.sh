#!/bin/bash

# Setup script for Git hooks

set -e

echo "🔧 Setting up Git hooks..."

# Copy pre-commit hook
cp scripts/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

echo "✅ Git hooks installed successfully!"
echo ""
echo "The following hooks are now active:"
echo "  • pre-commit: Runs Go lint, vet, and build checks"
echo ""
echo "To bypass hooks temporarily, use: git commit --no-verify"
