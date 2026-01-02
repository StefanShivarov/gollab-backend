#!/bin/bash

set -e
echo "Running unit tests..."
go test ./... -v -cover
