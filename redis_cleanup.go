package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func removeRedisChatStoreKeys() error {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL env var not set")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL env var not set")
	}

	rOpts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("failed to parse REDIS_URL: %v", err)
	}
	rdb := redis.NewClient(rOpts)
	defer rdb.Close()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return log.Output(2, "failed to create pool: "+err.Error())
	}
	defer pool.Close()

	rows, err := pool.Query(ctx, "SELECT store_id FROM tenant_channel WHERE channel_id = 64 and expired_at::date > '2026-04-01'")
	if err != nil {
		return log.Output(2, "failed to query tenant_channel: "+err.Error())
	}

	var storeIDs []string
	for rows.Next() {
		var storeID string
		if err := rows.Scan(&storeID); err != nil {
			log.Printf("row scan error: %v", err)
			continue
		}
		storeIDs = append(storeIDs, storeID)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return log.Output(2, "rows error: "+err.Error())
	}
	log.Printf("Fetched %d store IDs from DB (channel_id=64)", len(storeIDs))

	deleted := 0
	failed := 0

	for _, storeID := range storeIDs {
		key := "CHAT_STORE:" + storeID

		n, err := rdb.Del(ctx, key).Result()
		if err != nil {
			log.Printf("✗ %s delete failed: %v", key, err)
			failed++
			continue
		}

		if n > 0 {
			log.Printf("✓ %s deleted", key)
			deleted++
		} else {
			log.Printf("  %s not found (already absent)", key)
		}
	}

	log.Println("========== REDIS CLEANUP REPORT ==========")
	log.Printf("Total stores:  %d", len(storeIDs))
	log.Printf("Keys deleted:  %d", deleted)
	log.Printf("Keys failed:   %d", failed)
	log.Println("==========================================")

	return nil
}
