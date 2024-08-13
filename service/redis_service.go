package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var (
	ctx          = context.Background()
	rdb          *redis.Client
	cacheData    []redis.Z
	cacheExpiry  time.Time
	cacheTimeout = 5 * time.Second // Cache timeout duration
)

// InitializeRedisClient initializes the Redis client
func InitializeRedisClient() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr, // Use default DB
		Password: redisPassword,
	})
}

// SetUserGamePoints stores the user's game points in a Redis sorted set
func SetUserGamePoints(username string, points float64) (redis.Z, error) {
	// Check if the user already exists
	exists, err := rdb.ZScore(ctx, "user:points", username).Result()

	if err == redis.Nil {
		// User does not exist, so set their points to 0 first
		z := redis.Z{
			Score:  0,
			Member: username,
		}

		err := rdb.ZAdd(ctx, "user:points", &z).Err()
		if err != nil {
			return redis.Z{}, fmt.Errorf("could not set initial points for user %s: %v", username, err)
		}
		fmt.Printf("User %s did not exist. Set initial points to 0.\n", username)
	} else if err != nil {
		// An actual error occurred while checking the user's score
		return redis.Z{}, fmt.Errorf("could not check if user %s exists: %v", username, err)
	} else {
		// User exists, so return without setting points
		fmt.Printf("User %s already exists with score %f, not setting new points.\n", username, exists)
		return redis.Z{Score: exists, Member: username}, nil
	}

	// Return the newly added or existing user's score
	z := redis.Z{
		Score:  points,
		Member: username,
	}
	return z, nil
}

// IncrementUserGamePoints increments the user's points by a given value
func IncrementUserGamePoints(username string, increment int) (redis.Z, error) {
	newScore, err := rdb.ZIncrBy(ctx, "user:points", float64(increment), username).Result()
	if err != nil {
		return redis.Z{}, fmt.Errorf("could not increment game points for user %s: %v", username, err)
	}

	z := redis.Z{
		Score:  newScore,
		Member: username,
	}

	return z, nil
}

// GetAllUserPointsDesc retrieves all user points in descending order
func GetAllUserPointsDesc() []redis.Z {
	// Check if cache is valid
	if time.Now().Before(cacheExpiry) && cacheData != nil {
		fmt.Println("Returning cached data")
		return cacheData
	}

	// If cache is expired or empty, query Redis
	fmt.Println("Cache expired or empty, querying Redis")
	users, err := rdb.ZRevRangeWithScores(ctx, "user:points", 0, -1).Result()
	if err != nil {
		log.Fatalf("Could not retrieve user points: %v", err)
	}

	// Update cache
	cacheData = users
	cacheExpiry = time.Now().Add(cacheTimeout)

	return users
}
