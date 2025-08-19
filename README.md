# Chitauri Solana Balance API

This project is a Golang REST API (using Fiber) that fetches the Solana balance for one or more wallet addresses. It includes:
- IP rate limiting
- Per-wallet caching (10s TTL)
- Per-wallet mutex for concurrent requests
- MongoDB authentication (API keys)
- Dockerized setup with MongoDB

## Authentication

**Current:**
- The API uses simple static API tokens for authentication.
- Valid API keys are stored in the `api_keys` collection in MongoDB.
- Clients must include their API key in the `X-API-Key` header for every request.

**Planned:**
- In future versions, we will implement JWT authentication with access and refresh tokens for more robust, user-based security.
- JWTs will allow for stateless authentication, token expiry, and refresh flows.

## Usage

1. Start the API and MongoDB using Docker Compose:
   ```bash
   docker-compose up --build
   ```
2. Insert at least one API key into MongoDB:
   ```json
   { "key": "your-test-api-key" }
   ```
3. Make requests to `/api/get-balance` with the `X-API-Key` header and a JSON body:
   ```json
   {
     "wallets": ["WALLET_ADDRESS_1", "WALLET_ADDRESS_2"]
   }
   ```

## Environment Variables
- `SOLANA_API_KEY`: Your Helius Solana RPC API key
- `MONGODB_URI`: MongoDB connection string (default: `mongodb://mongo:27017`)
- `WALLET_ADDRESSES`: (optional) Comma-separated wallet addresses for testing

## Roadmap
- [x] API key authentication
- [ ] JWT authentication with access/refresh tokens
- [ ] User registration/login endpoints
- [ ] Token revocation and session management

---
For questions or contributions, open an issue or pull request.
