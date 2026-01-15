#!/bin/bash

# This script fixes image upload permissions for the Docker deployment
# Run this on your Ubuntu server

echo "Restoring permissions for the uploads directory..."

# Get the path of the current directory
PROJECT_DIR=$(pwd)
UPLOADS_DIR="$PROJECT_DIR/frontend/uploads"

# Check if the directory exists
if [ ! -d "$UPLOADS_DIR" ]; then
    echo "Creating uploads directory..."
    mkdir -p "$UPLOADS_DIR"
fi

# Set ownership to UID 1000 (standard Docker app user) and give write permissions
echo "Setting permissions on $UPLOADS_DIR"
chmod -R 777 "$UPLOADS_DIR"

# Also ensure the parent directories are accessible
chmod 755 "$PROJECT_DIR/frontend"

echo "✓ Upload permissions fixed!"
echo "Try uploading an image again. No restart needed."
