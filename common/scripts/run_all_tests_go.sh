#!/bin/bash

echo "Running tests for File Repository service..."

go test ./services/file-repository/... -v --timeout 30s

echo "Running tests for File Repository service: DONE"
