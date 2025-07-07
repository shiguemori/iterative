#!/bin/bash

# B3 Challenge API Benchmark Script
# This script tests the performance of the API endpoints

API_BASE_URL="http://localhost:8080"
CONCURRENT_REQUESTS=10
TOTAL_REQUESTS=1000

echo "=== B3 Challenge API Performance Benchmark ==="
echo "API Base URL: $API_BASE_URL"
echo "Concurrent Requests: $CONCURRENT_REQUESTS"
echo "Total Requests: $TOTAL_REQUESTS"
echo ""

# Check if API is running
echo "Checking if API is running..."
if ! curl -s "$API_BASE_URL/health" > /dev/null; then
    echo "Error: API is not running at $API_BASE_URL"
    echo "Please start the API first: go run cmd/api/main.go"
    exit 1
fi
echo "✓ API is running"
echo ""

# Install Apache Bench if not available
if ! command -v ab &> /dev/null; then
    echo "Installing Apache Bench (ab)..."
    sudo apt-get update && sudo apt-get install -y apache2-utils
fi

echo "=== Benchmark Results ==="
echo ""

# Test 1: Health Check Endpoint
echo "1. Health Check Endpoint (/health)"
echo "   Testing basic endpoint performance..."
ab -n $TOTAL_REQUESTS -c $CONCURRENT_REQUESTS "$API_BASE_URL/health" | grep -E "(Requests per second|Time per request|Transfer rate)"
echo ""

# Test 2: Get All Tickers
echo "2. Get All Tickers (/api/v1/tickers)"
echo "   Testing simple database query performance..."
ab -n 500 -c 5 "$API_BASE_URL/api/v1/tickers" | grep -E "(Requests per second|Time per request|Transfer rate)"
echo ""

# Test 3: Get Ticker Aggregation
echo "3. Get Ticker Aggregation (/api/v1/tickers/PETR4/aggregation)"
echo "   Testing complex aggregation query performance..."
ab -n 200 -c 5 "$API_BASE_URL/api/v1/tickers/PETR4/aggregation" | grep -E "(Requests per second|Time per request|Transfer rate)"
echo ""

# Test 4: Get Trades by Date Range
echo "4. Get Trades by Date Range (/api/v1/trades)"
echo "   Testing large dataset query performance..."
ab -n 100 -c 3 "$API_BASE_URL/api/v1/trades" | grep -E "(Requests per second|Time per request|Transfer rate)"
echo ""

echo "=== Benchmark Complete ==="
echo ""
echo "Performance Summary:"
echo "- Health check should handle >1000 req/sec"
echo "- Simple queries should handle >100 req/sec"
echo "- Complex aggregations should handle >10 req/sec"
echo "- Large dataset queries should handle >5 req/sec"
echo ""
echo "Note: Actual performance depends on hardware, database size, and system load."

