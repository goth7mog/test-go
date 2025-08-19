#!/bin/bash
# Test API caching behavior with precise timing and time deltas

LOGFILE="cache-test.log"
echo "--- Cache Test started at $(gdate +%s%3N) ---" > "$LOGFILE"

last_ts=0
log_request() {
  local name="$1"
  local cmd="$2"
  local ts=$(gdate +%s%3N)
  local result=$(eval "$cmd")
  local delta="N/A"
  if [ "$last_ts" -ne 0 ]; then
    delta=$((ts - last_ts))
  fi
  last_ts=$ts
  echo "$name - $ts - delta: $delta ms - $result" >> "$LOGFILE"
}

# Start
log_request "Request 1st" "curl -s -X POST http://localhost:8080/api/get-balance -H 'Content-Type: application/json' -H 'X-API-Key: d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5' -d '{\"wallets\":[\"9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b\"]}'"


log_request "Request 2nd" "curl -s -X POST http://localhost:8080/api/get-balance -H 'Content-Type: application/json' -H 'X-API-Key: d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5' -d '{\"wallets\":[\"9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b\"]}'"


log_request "Request 3rd" "curl -s -X POST http://localhost:8080/api/get-balance -H 'Content-Type: application/json' -H 'X-API-Key: d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5' -d '{\"wallets\":[\"9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b\"]}'"

# Wait for cache expiration
sleep 12

log_request "Request 4th (after 12s)" "curl -s -X POST http://localhost:8080/api/get-balance -H 'Content-Type: application/json' -H 'X-API-Key: d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5' -d '{\"wallets\":[\"9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b\"]}'"


log_request "Request 5th" "curl -s -X POST http://localhost:8080/api/get-balance -H 'Content-Type: application/json' -H 'X-API-Key: d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5' -d '{\"wallets\":[\"9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b\"]}'"


log_request "Request 6th" "curl -s -X POST http://localhost:8080/api/get-balance -H 'Content-Type: application/json' -H 'X-API-Key: d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5' -d '{\"wallets\":[\"9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b\"]}'"


log_request "Request 7th" "curl -s -X POST http://localhost:8080/api/get-balance -H 'Content-Type: application/json' -H 'X-API-Key: d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5' -d '{\"wallets\":[\"9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b\"]}'"

echo "--- Cache Test finished at $(gdate +%s%3N) ---" >> "$LOGFILE"
