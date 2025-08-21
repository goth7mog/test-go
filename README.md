# Chitauri Solana Balance API

This project is a Golang REST API (using Fiber) that fetches the Solana balance for one or more wallet addresses. It includes:
- IP rate limiting
- Per-wallet caching (10s TTL)
- Per-wallet mutex for concurrent requests
- MongoDB authentication (API keys)
- Dockerized setup with MongoDB

## Authentication
For the robust authentication it's recommended to implement JWT with access/refresh tokens going onwards

## Updates
### main.go
Added goroutines inside the `getBalanceHandler` function to fetch the wallet balances in parallel
### Nginx
Added for IP rate limiting
### Github Actions
Added

## Testing

### 1. Test with one wallet
```bash
curl -X POST http://localhost:8080/api/get-balance \
	-H "Content-Type: application/json" \
	-H "X-API-Key: d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5" \
	-d '{"wallets":["9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b"]}'
```

### 2. Test with multiple wallets
```bash
curl -X POST http://localhost:8080/api/get-balance \
	-H "Content-Type: application/json" \
	-H "X-API-Key: a9b8c7d6e5f4a3b2c1d0e9f8a7b6c5d4" \
	-d '{"wallets":["9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b","7BzJk1vQw2Qw1vQw2Qw1vQw2Qw1vQw2Qw1vQw2Qw2Qw2"]}'
```

### 3. Test with 5 requests using the same wallet (concurrent)
```bash
for i in {1..5}; do
	curl -X POST http://localhost:8080/api/get-balance \
		-H "Content-Type: application/json" \
		-H "X-API-Key: f1e2d3c4b5a6f7e8d9c0b1a2e3f4d5c6" \
		-d '{"wallets":["9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b"]}' &
done
wait
```

### 4. Test all the above at the same time
```bash
	./simultaneous-requests-test.sh
```

### 5. Test with invalid API key
```bash
curl -X POST http://localhost:8080/api/get-balance \
	-H "Content-Type: application/json" \
	-H "X-API-Key: invalid-key" \
	-d '{"wallets":["9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b"]}'
```

### 6. Test rate limiting (send >10 requests in a minute)
```bash
for i in {1..20}; do
	curl -X POST http://localhost:8080/api/get-balance \
		-H "Content-Type: application/json" \
		-H "X-API-Key: b2a3c4d5e6f7a8b9c0d1e2f3a4b5c6d7" \
		-d '{"wallets":["9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b"]}'
done
```

### 7. Cache
```bash
	./caching_test
```

