#!/bin/bash

# Stops execution when there was an error.
set -e

# Start the actual server application!
go run cmd/game-service/main.go