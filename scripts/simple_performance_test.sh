#!/bin/bash

# Simple Performance Test for B3 Challenge API
echo "=== B3 Challenge API Performance Test ==="
echo ""

API_BASE_URL="http://localhost:8080"

# Function to measure response time
measure_response_time() {
    local url=$1
    local description=$2
    echo "Testing: $description"
    echo "URL: $url"
    
    # Measure response time using curl
    local response_time=$(curl -o /dev/null -s -w "%{time_total}" "$url")
    local http_code=$(curl -o /dev/null -s -w "%{http_code}" "$url")
    
    echo "HTTP Status: $http_code"
    echo "Response Time: ${response_time}s"
    echo ""
}

# Test endpoints
measure_response_time "$API_BASE_URL/health" "Health Check"
measure_response_time "$API_BASE_URL/api/v1/tickers" "Get All Tickers"
measure_response_time "$API_BASE_URL/api/v1/tickers/PETR4/aggregation" "Get PETR4 Aggregation"
measure_response_time "$API_BASE_URL/api/v1/trades" "Get Last 7 Days Trades"

echo "=== Performance Test Complete ==="
echo ""
echo "Expected Performance Benchmarks:"
echo "- Health check: < 0.01s"
echo "- Simple queries: < 0.1s"
echo "- Complex aggregations: < 0.5s"
echo "- Large dataset queries: < 1.0s"

