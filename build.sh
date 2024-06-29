#!/bin/bash

# Create output directories if they do not exist
mkdir -p bin/plugins

# Clean up previous builds
echo "Cleaning up previous builds..."
rm -rf bin/plugins/*

# Build main application
echo "Building main application..."
cd src
go build -o ../bin/main main.go
cd ..

# Build plugins
# echo "Building plugins..."

# for plugin in plugins/*; do
#     if [ -d "$plugin" ]; then
#         plugin_name=$(basename "$plugin")
#         echo "Building $plugin_name..."
#         go build -buildmode=plugin -o bin/plugins/${plugin_name}.so ./$plugin/main.go
#     fi
# done

echo "Build complete."
