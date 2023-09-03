#!/bin/sh

response=$(echo "Hello" | nc $SERVER_IP $SERVER_PORT)

if [ "$response" == "Hello" ]; then
    echo "Server responded by repeating the message"
else
    echo "Server not responding"
fi