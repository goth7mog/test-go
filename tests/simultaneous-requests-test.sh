#!/bin/bash
# test_single_and_multiple_wallets.sh
# Simultaneously test the API with 1 wallet and multiple wallets

# Log file
LOGFILE="simultaneous-requests.log"
START_TS=$(gdate +%s%3N)
echo "--- Test started at $START_TS ---" > "$LOGFILE"

log_request() {
  local name="$1"
  local cmd="$2"
  (
    local ts=$(gdate +%s%3N)
    local result=$(eval "$cmd")
    local delta=$((ts - START_TS))
    echo "$name - $ts (delta ${delta}ms) - $result" >> "$LOGFILE"
  ) &
}

# Single wallet request
log_request "Single Wallet request" "curl -s -X POST http://localhost:8080/api/get-balance -H 'Content-Type: application/json' -H 'X-API-Key: d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5' -d '{\"wallets\":[\"9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b\"]}'"

# Multiple wallets request
log_request "Multiple Wallets request" "curl -s -X POST http://localhost:8080/api/get-balance -H 'Content-Type: application/json' -H 'X-API-Key: a9b8c7d6e5f4a3b2c1d0e9f8a7b6c5d4' -d '{\"wallets\":[\"9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b\",\"7BzJk1vQw2Qw1vQw2Qw1vQw2Qw1vQw2Qw1vQw2Qw2Qw2\"]}'"

# 5 concurrent requests for the same wallet
for i in {1..5}; do
  log_request "Concurrent request $i" "curl -s -X POST http://localhost:8080/api/get-balance -H 'Content-Type: application/json' -H 'X-API-Key: f1e2d3c4b5a6f7e8d9c0b1a2e3f4d5c6' -d '{\"wallets\":[\"9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b\"]}'"
done

wait
END_TS=$(gdate +%s%3N)
TOTAL_DELTA=$((END_TS - START_TS))
echo "--- Test finished at $END_TS (total duration ${TOTAL_DELTA}ms) ---" >> "$LOGFILE"
