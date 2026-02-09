#!/bin/bash

# Chicago Poker - Server Start Script
# Starts the server and provides instructions for connecting with netcat

echo "==================================="
echo "   Chicago Poker Server"
echo "==================================="

# Kill any existing server
pkill -f chicago-poker 2>/dev/null
sleep 1

# Build if needed
if [ ! -f "./chicago-poker" ]; then
    echo "Building game..."
    export PATH=$PATH:/usr/local/go/bin
    go build ./cmd/chicago-poker || { echo "Build failed!"; exit 1; }
fi

# Start server in background
export PATH=$PATH:/usr/local/go/bin
echo "Starting server on port 8080..."
./chicago-poker 2>&1 | tee server.log &
SERVER_PID=$!
sleep 1

if ! ps -p $SERVER_PID > /dev/null; then
    echo "ERROR: Server failed to start. Check server.log for details."
    exit 1
fi

echo ""
echo "✓ Server started (PID: $SERVER_PID)"
echo ""
echo "==================================="
echo "HOW TO PLAY:"
echo "==================================="
echo ""
echo "1. Open 2 new terminals and connect:"
echo "   nc localhost 8080"
echo ""
echo "2. Enter usernames when prompted"
echo ""
echo "3. Game starts automatically with 2 players"
echo ""
echo "4. Gameplay:"
echo "   • Poker Round: Type card indices to toss (e.g., '0 2' or press Enter to keep all)"
echo "   • Trick Round: Type a single card index (e.g., '0')"
echo "   • First to 50 points wins!"
echo ""
echo "==================================="
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

# Wait for Ctrl+C
trap "kill $SERVER_PID 2>/dev/null; pkill -f chicago-poker 2>/dev/null; echo ''; echo '✓ Server stopped'; rm -f server.log; exit 0" INT TERM

wait
