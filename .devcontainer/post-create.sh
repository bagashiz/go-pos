#!/bin/sh

# download dependencies
echo "Downloading dependencies..."
go mod download

# install task runner
echo "Installing task runner..."
go install github.com/go-task/task/v3/cmd/task@latest

# install dependencies, migrate database, and run server
echo "Installing dependencies, migrating database, and running server..."
task 