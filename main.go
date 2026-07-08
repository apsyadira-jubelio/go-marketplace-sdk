package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/lazada"
	"github.com/apsyadira-jubelio/go-marketplace-sdk/shopee"
	"github.com/davecgh/go-spew/spew"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type ShopeeChannelExtraInfo struct {
	Mode            string      `json:"mode"`
	Secret          string      `json:"secret"`
	ShopID          uint64      `json:"shopId"`
	Expired         bool        `json:"expired"`
	Version         json.Number `json:"version"`
	PartnerID       int         `json:"partnerId"`
	AccessToken     string      `json:"access_token"`
	RefreshToken    string      `json:"refresh_token"`
	TokenExpiredAt  string      `json:"token_expired_at"`
	NewRefreshToken string      `json:"new_refresh_token"`
}

func main() {
	godotenv.Load()

	cmd := "refresh"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "refresh":
		if err := refreshShopeeTokensFromDB(); err != nil {
			log.Fatal(err)
		}
	case "redis-cleanup":
		if err := removeRedisChatStoreKeys(); err != nil {
			log.Fatal(err)
		}
	case "conversations":
		if err := getListConversation(); err != nil {
			log.Fatal(err)
		}
	case "one-conversation":
		if err := getOneConversation(); err != nil {
			log.Fatal(err)
		}
	case "lazada-conversations":
		if err := getLazadaConversations(); err != nil {
			log.Fatal(err)
		}
	case "lazada-concurrent":
		if err := lazadaConcurrentExample(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown command: %s (use: refresh | redis-cleanup | conversations | one-conversation | lazada-conversations)", cmd)
	}
}

func getListConversation() error {
	partnerID, _ := strconv.Atoi(os.Getenv("SHOPEE_PARTNER_ID"))
	appConfig := shopee.AppConfig{
		APIURL:      os.Getenv("SHOPEE_API_URL"),
		PartnerID:   partnerID,
		PartnerKey:  os.Getenv("SHOPEE_PARTNER_KEY"),
		RedirectURL: os.Getenv("SHOPEE_REDIRECT_URL"),
	}

	var opts []shopee.Option
	if socksAddr := os.Getenv("SHOPEE_SOCKS5_ADDR"); socksAddr != "" {
		opts = append(opts, shopee.WithSocks5(socksAddr))
	}
	shopeeClient := shopee.NewClient(appConfig, opts...)

	shopID, _ := strconv.ParseUint(os.Getenv("SHOPEE_SHOP_ID"), 10, 64)
	accessToken := os.Getenv("SHOPEE_ACCESS_TOKEN")

	resp, err := shopeeClient.Chat.GetConversationList(shopID, accessToken, shopee.GetConversationParamsRequest{
		Direction: "latest",
		Type:      "all",
		PageSize:  20,
	})

	if err != nil {
		return fmt.Errorf("get conversation list: %w", err)
	}

	log.Printf("Conversations: %d", len(resp.Response.ConversationsList))
	spew.Dump(resp)

	return nil
}

func getOneConversation() error {
	partnerID, _ := strconv.Atoi(os.Getenv("SHOPEE_PARTNER_ID"))
	appConfig := shopee.AppConfig{
		APIURL:      os.Getenv("SHOPEE_API_URL"),
		PartnerID:   partnerID,
		PartnerKey:  os.Getenv("SHOPEE_PARTNER_KEY"),
		RedirectURL: os.Getenv("SHOPEE_REDIRECT_URL"),
	}

	var opts []shopee.Option
	if socksAddr := os.Getenv("SHOPEE_SOCKS5_ADDR"); socksAddr != "" {
		opts = append(opts, shopee.WithSocks5(socksAddr))
	}
	opts = append(opts, shopee.WithLogger(&shopee.LeveledLogger{Level: shopee.LevelDebug}))
	shopeeClient := shopee.NewClient(appConfig, opts...)

	shopID, _ := strconv.ParseUint(os.Getenv("SHOPEE_SHOP_ID"), 10, 64)
	accessToken := os.Getenv("SHOPEE_ACCESS_TOKEN")
	conversationID := os.Getenv("SHOPEE_CONVERSATION_ID")

	log.Printf("shopID=%d conversationID=%s accessTokenLen=%d", shopID, conversationID, len(accessToken))

	// Paginate through all conversations to find the target conversation_id
	targetConversationID := conversationID
	var nextTimeNano int64
	pageSize := 50
	found := false
	totalScanned := 0

	for {
		params := shopee.GetConversationParamsRequest{
			Direction: "latest",
			Type:      "all",
			PageSize:  pageSize,
		}
		if nextTimeNano > 0 {
			params.NextTimeNano = nextTimeNano
		}

		listResp, err := shopeeClient.Chat.GetConversationList(shopID, accessToken, params)
		if err != nil {
			return fmt.Errorf("get conversation list page: %w", err)
		}

		for _, c := range listResp.Response.ConversationsList {
			totalScanned++
			if c.ConversationID == targetConversationID {
				log.Printf("FOUND! conversation_id=%s to_name=%s", c.ConversationID, c.ToName)
				found = true
			}
		}

		log.Printf("Scanned %d conversations (page total: %d, more: %v)", totalScanned, len(listResp.Response.ConversationsList), listResp.Response.PageResult.More)

		if !listResp.Response.PageResult.More || len(listResp.Response.ConversationsList) == 0 {
			break
		}

		// Parse next cursor timestamp
		if ts := listResp.Response.PageResult.NextCursor.NextMessageTimeNano; ts != "" && ts != "0" {
			nextTimeNano, _ = strconv.ParseInt(ts, 10, 64)
		} else {
			break
		}
	}

	log.Printf("Total scanned: %d, found: %v", totalScanned, found)
	if !found {
		log.Printf("conversation_id=%s not found in any page for shop_id=%d", targetConversationID, shopID)
		return nil
	}

	conversationIDInt, err := strconv.ParseInt(conversationID, 10, 64)
	if err != nil {
		return fmt.Errorf("parse conversation_id: %w", err)
	}

	resp, err := shopeeClient.Chat.GetOneConversation(shopID, accessToken, shopee.GetMessageParamsRequest{
		ConversationID: conversationIDInt,
	})

	if err != nil {
		return fmt.Errorf("get one conversation: %w", err)
	}

	log.Printf("Conversation detail:")
	spew.Dump(resp)

	return nil
}

func getLazadaConversations() error {
	ctx := context.Background()

	appKey := os.Getenv("LAZADA_APP_KEY")
	secret := os.Getenv("LAZADA_SECRET_KEY")
	accessToken := os.Getenv("LAZADA_ACCESS_TOKEN")
	region := lazada.Region(os.Getenv("LAZADA_REGION"))
	if region == "" {
		region = lazada.Indonesia
	}

	lazadaClient := lazada.NewClient(appKey, secret, region)

	// Token passed explicitly per-request (race-free for concurrent usage)
	resp, err := lazadaClient.Chat.GetSessionList(ctx, accessToken, &lazada.SessionListQuery{
		PageSize:  20,
		StartTime: time.Now().AddDate(0, -1, 0).UnixNano() / int64(time.Millisecond),
	})
	if err != nil {
		return fmt.Errorf("get lazada session list: %w", err)
	}

	log.Printf("Lazada sessions: %d", len(resp.SessionListResponseData.SessionList))
	spew.Dump(resp)

	return nil
}

// lazadaConcurrentExample demonstrates race-free concurrent requests with per-tenant tokens.
// Pass token explicitly to each service method instead of NewTokenClient().
// NewTokenClient() mutates client state and causes race conditions when used concurrently.
func lazadaConcurrentExample() error {
	ctx := context.Background()

	appKey := os.Getenv("LAZADA_APP_KEY")
	secret := os.Getenv("LAZADA_SECRET_KEY")
	region := lazada.Region(os.Getenv("LAZADA_REGION"))
	if region == "" {
		region = lazada.Indonesia
	}

	// Single client instance shared across goroutines (no token on client)
	client := lazada.NewClient(appKey, secret, region)

	// Simulate multiple tenants with different tokens
	tenants := []struct {
		tenantID string
		token    string
	}{
		{"tenant-1", os.Getenv("LAZADA_ACCESS_TOKEN")},
		{"tenant-2", os.Getenv("LAZADA_ACCESS_TOKEN_2")},
		{"tenant-3", os.Getenv("LAZADA_ACCESS_TOKEN_3")},
	}

	results := make(chan struct {
		tenantID string
		count    int
		err      error
	}, len(tenants))

	for _, tenant := range tenants {
		go func(t struct {
			tenantID string
			token    string
		}) {
			// Pass tenant-specific token explicitly (race-free)
			resp, err := client.Chat.GetSessionList(ctx, t.token, &lazada.SessionListQuery{
				PageSize:  20,
				StartTime: time.Now().AddDate(0, -1, 0).UnixNano() / int64(time.Millisecond),
			})
			if err != nil {
				results <- struct {
					tenantID string
					count    int
					err      error
				}{t.tenantID, 0, err}
				return
			}
			results <- struct {
				tenantID string
				count    int
				err      error
			}{t.tenantID, len(resp.SessionListResponseData.SessionList), nil}
		}(tenant)
	}

	var errors []string
	for i := 0; i < len(tenants); i++ {
		result := <-results
		if result.err != nil {
			log.Printf("[%s] ERROR: %v", result.tenantID, result.err)
			errors = append(errors, result.tenantID)
		} else {
			log.Printf("[%s] OK: %d sessions", result.tenantID, result.count)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed tenants: %v", errors)
	}

	log.Println("All concurrent requests completed successfully")
	return nil
}

func refreshShopeeTokensFromDB() error {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL env var not set")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return log.Output(2, "failed to create pool: "+err.Error())
	}
	defer pool.Close()

	apiURL := os.Getenv("SHOPEE_API_URL")
	if apiURL == "" {
		apiURL = "https://partner.shopeemobile.com"
	}
	redirectURL := os.Getenv("SHOPEE_REDIRECT_URL")
	socksAddr := os.Getenv("SHOPEE_SOCKS5_ADDR")

	// Step 1: Fetch all rows into memory (releases DB cursor early)
	type channelRow struct {
		storeID   string
		tenantID  string
		storeName string
		info      ShopeeChannelExtraInfo
		isActive  bool
		expiredAt time.Time
	}
	var channels []channelRow

	rows, err := pool.Query(ctx, "SELECT store_id, tenant_id, store_name, is_active, expired_at, extra_info::text FROM tenant_channel WHERE channel_id = 64 AND expired_at >= '2026-07-03'")
	if err != nil {
		return log.Output(2, "failed to query tenant_channel: "+err.Error())
	}
	for rows.Next() {
		var storeID string
		var tenantID string
		var storeName string
		var isActive bool
		var expiredAt time.Time
		var extraInfoJSON string
		if err := rows.Scan(&storeID, &tenantID, &storeName, &isActive, &expiredAt, &extraInfoJSON); err != nil {
			log.Printf("row scan error: %v", err)
			continue
		}
		var info ShopeeChannelExtraInfo
		if err := json.Unmarshal([]byte(extraInfoJSON), &info); err != nil {
			log.Printf("[%s] failed to parse extra_info: %v", storeID, err)
			continue
		}
		channels = append(channels, channelRow{storeID: storeID, tenantID: tenantID, storeName: storeName, info: info, isActive: isActive, expiredAt: expiredAt})
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return log.Output(2, "rows error: "+err.Error())
	}
	log.Printf("Fetched %d channels from DB", len(channels))

	// Step 2: Refresh tokens (reuse client when partnerID+secret are the same)
	var lastPartnerID int
	var lastSecret string
	var client *shopee.ShopeeClient

	type updateRow struct {
		storeID   string
		storeName string
		info      ShopeeChannelExtraInfo
		expiredAt time.Time
	}
	var updates []updateRow
	var failures []string

	for _, ch := range channels {
		info := ch.info
		storeID := ch.storeID
		storeName := ch.storeName

		if client == nil || info.PartnerID != lastPartnerID || info.Secret != lastSecret {
			appConfig := shopee.AppConfig{
				PartnerID:   info.PartnerID,
				PartnerKey:  info.Secret,
				RedirectURL: redirectURL,
				APIURL:      apiURL,
			}
			if socksAddr != "" {
				client = shopee.NewClient(appConfig, shopee.WithSocks5(socksAddr))
			} else {
				client = shopee.NewClient(appConfig)
			}
			lastPartnerID = info.PartnerID
			lastSecret = info.Secret
		}

		respRefresh, err := client.Auth.RefreshAccessToken(info.ShopID, 0, info.RefreshToken)
		if err != nil && info.NewRefreshToken != "" {
			log.Printf("[%s] %s shopID=%d refresh_token failed, retrying with new_refresh_token: %v", storeID, storeName, info.ShopID, err)
			respRefresh, err = client.Auth.RefreshAccessToken(info.ShopID, 0, info.NewRefreshToken)
		}
		if err != nil {
			log.Printf("[%s] %s shopID=%d refresh failed: %v", storeID, storeName, info.ShopID, err)
			failures = append(failures, storeName)
			continue
		}

		log.Printf("[%s] %s shopID=%d token refreshed! ExpireIn: %d", storeID, storeName, info.ShopID, respRefresh.ExpireIn)

		info.AccessToken = respRefresh.AccessToken
		info.RefreshToken = respRefresh.RefreshToken
		info.NewRefreshToken = ""
		info.TokenExpiredAt = time.Now().UTC().Add(time.Duration(respRefresh.ExpireIn) * time.Second).Format(time.RFC3339)
		info.Expired = false

		updates = append(updates, updateRow{storeID: storeID, storeName: storeName, info: info, expiredAt: time.Now().UTC().Add(time.Duration(respRefresh.ExpireIn) * time.Second)})
	}
	log.Printf("%d/%d tokens refreshed successfully", len(updates), len(channels))

	// Step 3: Batch all DB updates in a single transaction
	if len(updates) == 0 {
		return nil
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return log.Output(2, "failed to begin transaction: "+err.Error())
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}
	for _, u := range updates {
		updatedJSON, err := json.Marshal(u.info)
		if err != nil {
			log.Printf("[%s] failed to marshal updated extra_info: %v", u.storeID, err)
			continue
		}
		batch.Queue("UPDATE tenant_channel SET extra_info = $1, is_active = true, expired_at = $2 WHERE store_id = $3", updatedJSON, u.expiredAt, u.storeID)
	}

	br := tx.SendBatch(ctx, batch)
	for range updates {
		if _, err := br.Exec(); err != nil {
			log.Printf("batch update error: %v", err)
		}
	}
	br.Close()

	if err := tx.Commit(ctx); err != nil {
		return log.Output(2, "failed to commit transaction: "+err.Error())
	}

	log.Printf("Batch update complete: %d rows updated in a single transaction", len(updates))

	// Step 4: Report
	log.Println("========== REPORT ==========")
	log.Printf("Total channels: %d", len(channels))
	log.Printf("Total sukses:   %d", len(updates))
	log.Printf("Total gagal:    %d", len(failures))
	if len(updates) > 0 {
		log.Println("--- Sukses ---")
		for _, u := range updates {
			log.Printf("  ✓ %s (shopID=%d)", u.storeName, u.info.ShopID)
		}
	}
	if len(failures) > 0 {
		log.Println("--- Gagal ---")
		for _, name := range failures {
			log.Printf("  ✗ %s", name)
		}
	}
	log.Println("============================")

	return nil
}
