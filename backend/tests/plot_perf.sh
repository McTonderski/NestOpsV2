#!/bin/bash

# Path to the Go application binary
GO_APP="./your_go_app_binary"

# Log file for response times
RESPONSE_TIMES_LOG="response_times.log"

# DDoSify binary path
DDOSIFY="ddosify"

# Target URL for DDoSify
TARGET_URL="http://localhost:8080"

# Start the Go application in the background
echo "Starting Go application..."
$GO_APP &
GO_APP_PID=$!

# Wait for the Go app to start
sleep 5

# Clear the response times log
> $RESPONSE_TIMES_LOG

# Incremental request count
REQUEST_COUNT=10
MAX_REQUESTS=100

echo "Starting DDoSify with incremental request counts..."
while [ $REQUEST_COUNT -le $MAX_REQUESTS ]; do
    echo "Running DDoSify with $REQUEST_COUNT requests..."
    
    # Run DDoSify and capture response times
    $DDOSIFY -t $TARGET_URL -n $REQUEST_COUNT | tee -a $RESPONSE_TIMES_LOG
    
    # Increment the request count
    REQUEST_COUNT=$((REQUEST_COUNT + 10))
    
    # Wait a bit before the next run
    sleep 2
done

# Stop the Go application
echo "Stopping Go application..."
kill $GO_APP_PID

echo "Test completed. Response times logged in $RESPONSE_TIMES_LOG."