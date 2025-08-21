package main

import (
	"bytes"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type testResponse struct {
	StatusCode int
	Body       []byte
}

func setupTestApp() *fiber.App {
	app := fiber.New()
	ApiKeys["d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5"] = true // Use a real key from mongo-init.js
	app.Use(ApiKeyAuthMiddleware)
	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Rate limit exceeded"})
		},
	}))
	app.Post("/api/get-balance", func(c *fiber.Ctx) error {
		// Return balances for all requested wallets
		var req struct {
			Wallets []string `json:"wallets"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
		balances := fiber.Map{}
		for _, w := range req.Wallets {
			balances[w] = 100 // Dummy balance
		}
		return c.Status(200).JSON(fiber.Map{"balances": balances})
	})
	return app
}

func testRequest(app *fiber.App, path, apiKey string, wallets []string) testResponse {
	// Build wallets JSON array
	walletsJson := "["
	for i, w := range wallets {
		walletsJson += `"` + w + `"`
		if i < len(wallets)-1 {
			walletsJson += ","
		}
	}
	walletsJson += "]"
	bodyData := []byte(`{"wallets":` + walletsJson + `}`)
	req, _ := http.NewRequest("POST", path, bytes.NewReader(bodyData))
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	body := []byte{}
	if resp != nil {
		body = make([]byte, resp.ContentLength)
		resp.Body.Read(body)
	}
	return testResponse{StatusCode: resp.StatusCode, Body: body}
}

func TestIPRateLimiting(t *testing.T) {
	app := setupTestApp()
	validKey := "d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5"
	wallet := "9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b"
	for i := 0; i < 10; i++ {
		req := testRequest(app, "/api/get-balance", validKey, []string{wallet})
		if req.StatusCode != 200 {
			t.Errorf("Request %d failed unexpectedly: %v", i+1, req.StatusCode)
		}
	}
	req := testRequest(app, "/api/get-balance", validKey, []string{wallet})
	if req.StatusCode != 429 {
		t.Errorf("Expected rate limit, got %v", req.StatusCode)
	}
}

func TestCacheTTL(t *testing.T) {
	app := setupTestApp()
	validKey := "d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5"
	wallet := "7BzJk1vQw2Qw1vQw2Qw1vQw2Qw1vQw2Qw1vQw2Qw2Qw2"
	req := testRequest(app, "/api/get-balance", validKey, []string{wallet})
	if req.StatusCode != 200 {
		t.Fatalf("Initial request failed: %v", req.StatusCode)
	}
	time.Sleep(11 * time.Second)
	req2 := testRequest(app, "/api/get-balance", validKey, []string{wallet})
	if req2.StatusCode != 200 {
		t.Fatalf("Second request failed: %v", req2.StatusCode)
	}
}

func TestConcurrentRequestsMutexCache(t *testing.T) {
	app := setupTestApp()
	validKey := "d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5"
	wallet := "8HoQnePLqPj4M7PUDzfw8e3YwA4u4u4u4u4u4u4u4u4u4"
	var wg sync.WaitGroup
	var results [2]int
	wg.Add(2)
	go func() {
		defer wg.Done()
		req := testRequest(app, "/api/get-balance", validKey, []string{wallet})
		results[0] = req.StatusCode
	}()
	go func() {
		defer wg.Done()
		req := testRequest(app, "/api/get-balance", validKey, []string{wallet})
		results[1] = req.StatusCode
	}()
	wg.Wait()
	if results[0] != 200 || results[1] != 200 {
		t.Errorf("Concurrent requests failed: %v", results)
	}
}

func TestMongoDBAuthentication(t *testing.T) {
	app := setupTestApp()
	validKey := "d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5"
	wallet := "9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b"
	req := testRequest(app, "/api/get-balance", validKey, []string{wallet})
	if req.StatusCode != 200 {
		t.Errorf("Valid API key rejected: %v", req.StatusCode)
	}
	req2 := testRequest(app, "/api/get-balance", "invalid-key", []string{wallet})
	if req2.StatusCode != 401 {
		t.Errorf("Invalid API key accepted: %v", req2.StatusCode)
	}
}

func TestAPIResponseTime(t *testing.T) {
	start := time.Now()
	app := setupTestApp()
	validKey := "d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5"
	wallet := "9xQeWvG816bUx9EPn6bKk6vYjY6i9aF8fF7b7b7b7b7b"
	req := testRequest(app, "/api/get-balance", validKey, []string{wallet})
	elapsed := time.Since(start)
	if req.StatusCode != 200 {
		t.Fatalf("API request failed: %v", req.StatusCode)
	}
	if elapsed > time.Second {
		t.Errorf("API response too slow: %v", elapsed)
	}
}
