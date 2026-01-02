#!/bin/bash

set -e

echo "Downloading dependencies..."
go mod download

echo "Building Go application..."
mkdir -p bin
go build -buildvcs=false -o bin/myapp ./cmd
