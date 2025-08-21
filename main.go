package main

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Load wallet addresses from .env file
func loadWalletAddresses() []string {
	godotenv.Load()
	wallets := os.Getenv("WALLET_ADDRESSES")
	if wallets == "" {
		return []string{}
	}
	return strings.Split(wallets, ",")
}

var (
	rpcClient   *rpc.Client
	walletCache = cache.New(10*time.Second, 1*time.Minute)
	walletLocks = struct {
		sync.Mutex
		m map[string]*sync.Mutex
	}{m: make(map[string]*sync.Mutex)}
	apiKeys = make(map[string]bool)
)

func getWalletMutex(wallet string) *sync.Mutex {
	walletLocks.Lock()
	defer walletLocks.Unlock()
	if _, exists := walletLocks.m[wallet]; !exists {
		walletLocks.m[wallet] = &sync.Mutex{}
	}
	return walletLocks.m[wallet]
}

func loadAPIKeysFromMongo() {
	mongoUri := os.Getenv("MONGODB_URI")
	if mongoUri == "" {
		mongoUri = "mongodb://mongo:27017" // default for docker-compose
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}
	collection := client.Database("chitauri").Collection("api_keys")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
	for cursor.Next(context.Background()) {
		var result struct {
			Key string `bson:"key"`
		}
		if err := cursor.Decode(&result); err == nil {
			apiKeys[result.Key] = true
		}
	}
	cursor.Close(context.Background())
}

func apiKeyAuthMiddleware(c *fiber.Ctx) error {
	key := c.Get("X-API-Key")
	if !apiKeys[key] {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid API key"})
	}
	return c.Next()
}

type BalanceRequest struct {
	Wallets []string `json:"wallets"`
}

func getBalanceHandler(c *fiber.Ctx) error {
	var req BalanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	results := make(map[string]uint64)
	errors := make(map[string]map[string]string)
	var wg sync.WaitGroup
	var resultsMutex sync.Mutex
	for _, wallet := range req.Wallets {
		wg.Add(1)
		go func(wallet string) {
			defer wg.Done()
			mutex := getWalletMutex(wallet)
			mutex.Lock()
			defer mutex.Unlock()
			if cachedBalance, found := walletCache.Get(wallet); found {
				resultsMutex.Lock()
				results[wallet] = cachedBalance.(uint64)
				resultsMutex.Unlock()
				return
			}
			pubKey := solana.MustPublicKeyFromBase58(wallet)
			balance, err := rpcClient.GetBalance(context.Background(), pubKey, rpc.CommitmentFinalized)
			resultsMutex.Lock()
			if err != nil {
				results[wallet] = 0
				errors[wallet] = map[string]string{
					"code":    "RPC_ERROR",
					"message": "RPC call failed: " + err.Error(),
				}
			} else {
				results[wallet] = balance.Value
				walletCache.Set(wallet, balance.Value, cache.DefaultExpiration)
			}
			resultsMutex.Unlock()
		}(wallet)
	}
	wg.Wait()
	return c.JSON(fiber.Map{"balances": results, "errors": errors})
}

func main() {
	godotenv.Load()
	loadAPIKeysFromMongo()

	// Load Solana RPC API key and host from .env
	solanaApiKey := os.Getenv("SOLANA_API_KEY")
	solanaRpcHost := os.Getenv("SOLANA_RPC_HOST")
	rpcUrl := "https://" + solanaRpcHost + "/?api-key=" + solanaApiKey
	rpcClient = rpc.New(rpcUrl)

	app := fiber.New()

	// IP rate limiting: 10 requests per minute per IP
	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Rate limit exceeded"})
		},
	}))

	app.Use(apiKeyAuthMiddleware)
	app.Post("/api/get-balance", getBalanceHandler)
	app.Listen(":8080")
}
