#!/bin/bash

# Check if a port number is provided
if [ -z "$1" ]; then
  echo "Usage: $0 {PORT}"
  exit 1
fi

PORT=$1

# Get PIDs using lsof and kill them
PIDS=$(lsof -ti tcp:"$PORT")

if [ -z "$PIDS" ]; then
  echo "No processes found on TCP port $PORT."
else
  echo "Killing processes on TCP port $PORT: $PIDS"
  kill -9 $PIDS
  echo "Processes killed."
fi
