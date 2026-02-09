#!/bin/bash

# Automated test for Chicago Poker
# Tests basic connection and game start

echo "==================================="
echo "Chicago Poker - Automated Test"
echo "==================================="

# Kill any existing processes
pkill -f chicago-poker 2>/dev/null
pkill -f "nc localhost 8080" 2>/dev/null
sleep 1

export PATH=$PATH:/usr/local/go/bin

# Start server
echo "Starting server..."
./chicago-poker > test_server.log 2>&1 &
SERVER_PID=$!
sleep 2

# Check server started
if ! ps -p $SERVER_PID > /dev/null; then
    echo "❌ Server failed to start"
    cat test_server.log
    exit 1
fi
echo "✓ Server started"

# Connect player 1
echo "Connecting Player1..."
{
    sleep 1
    echo "Player1"
    sleep 2
    echo ""
    sleep 2
    echo "0"
    sleep 20
} | nc localhost 8080 > player1.log 2>&1 &
P1_PID=$!
sleep 1

# Connect player 2
echo "Connecting Player2..."
{
    sleep 1
    echo "Player2"
    sleep 2
    echo ""
    sleep 2
    echo "0"
    sleep 20
} | nc localhost 8080 > player2.log 2>&1 &
P2_PID=$!
sleep 5

# Check outputs
echo ""
echo "==================================="
echo "Test Results:"
echo "==================================="

if grep -q "Welcome to Chicago Poker" player1.log; then
    echo "✓ Player 1 received welcome message"
else
    echo "❌ Player 1 did not receive welcome message"
fi

if grep -q "Welcome to Chicago Poker" player2.log; then
    echo "✓ Player 2 received welcome message"
else
    echo "❌ Player 2 did not receive welcome message"
fi

if grep -q "Starting the game" player1.log || grep -q "Starting the game" player2.log; then
    echo "✓ Game started with 2 players"
else
    echo "❌ Game did not start"
fi

if grep -q "Your Hand" player1.log; then
    echo "✓ Player 1 received hand display"
else
    echo "❌ Player 1 did not receive hand display"
fi

if grep -q "Your Hand" player2.log; then
    echo "✓ Player 2 received hand display"
else
    echo "❌ Player 2 did not receive hand display"
fi

echo ""
echo "==================================="
echo "Player 1 Output (first 30 lines):"
echo "==================================="
head -30 player1.log

echo ""
echo "==================================="
echo "Player 2 Output (first 30 lines):"
echo "==================================="
head -30 player2.log

# Cleanup
kill $SERVER_PID $P1_PID $P2_PID 2>/dev/null
sleep 1
pkill -f chicago-poker 2>/dev/null
pkill -f "nc localhost 8080" 2>/dev/null

# Clean up log files
rm -f test_server.log player1.log player2.log 2>/dev/null

echo ""
echo "✓ Test complete and cleaned up"
