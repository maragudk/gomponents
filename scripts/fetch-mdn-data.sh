#!/bin/bash
# Fetch MDN browser-compat-data for HTML generation.
# This clones a shallow copy of the repository to minimize disk usage.

set -e

REPO_URL="https://github.com/mdn/browser-compat-data.git"
TARGET_DIR="tmp/browser-compat-data"

# Change to repository root
cd "$(dirname "$0")/.."

# Create tmp directory if needed
mkdir -p tmp

if [ -d "$TARGET_DIR" ]; then
    echo "Updating existing MDN data..."
    cd "$TARGET_DIR"
    git pull --depth 1
else
    echo "Cloning MDN browser-compat-data (shallow clone)..."
    git clone --depth 1 --filter=blob:none --sparse "$REPO_URL" "$TARGET_DIR"
    cd "$TARGET_DIR"
    # Only checkout the html directory we need
    git sparse-checkout set html
fi

echo "Done. MDN data available at $TARGET_DIR"
