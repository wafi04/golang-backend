package testing

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wafi04/golang-backend/grpc/pb"
	handler "github.com/wafi04/golang-backend/services/order/internal"
)

const (
	testVariantID = "cepet"
	initialStock  = 100
	totalOrders   = 110
	testRedisAddr = ""
	testDBConn    = ""
	testRedisPass = ""
	testRedisDB   = 0
)

func setupDatabaseAndRedis(t *testing.T) (*sqlx.DB, *redis.Client) {
	t.Helper()

	// Setup database
	db := sqlx.MustConnect("pgx", testDBConn)
	tx, err := db.Beginx()
	require.NoError(t, err, "Failed to start transaction")

	_, err = tx.Exec("DELETE FROM stock WHERE variant_id = $1", testVariantID)
	require.NoError(t, err, "Failed to clean up previous test data")

	_, err = tx.Exec(`
	INSERT INTO stock (variant_id, quantity, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW())`,
		testVariantID, initialStock,
	)
	require.NoError(t, err, "Failed to seed initial stock")

	err = tx.Commit()
	require.NoError(t, err, "Failed to commit transaction")

	// Setup Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     testRedisAddr,
		Password: testRedisPass,
		DB:       testRedisDB,
	})
	require.NoError(t, redisClient.FlushDB(context.Background()).Err(), "Failed to flush Redis DB")

	return db, redisClient
}

func tearDownDatabaseAndRedis(db *sqlx.DB, redisClient *redis.Client, t *testing.T) {
	t.Helper()

	// Cleanup database
	_, err := db.Exec("DELETE FROM orders; DELETE FROM stock WHERE variant_id = $1", testVariantID)
	require.NoError(t, err, "Failed to clean up database")

	// Cleanup Redis
	require.NoError(t, redisClient.FlushDB(context.Background()).Err(), "Failed to flush Redis DB")
}

func TestOrderService_ConcurrentOrderCreation(t *testing.T) {
	t.Parallel() // Allow parallel execution of tests

	// ======================
	// 1. setup dependencies
	// ======================
	db, redisClient := setupDatabaseAndRedis(t)
	defer tearDownDatabaseAndRedis(db, redisClient, t)

	orderService := handler.NewOrderService(db, redisClient)

	// ======================
	// 2. COncurency controlee
	// ======================
	var (
		wg        sync.WaitGroup
		successes int
		errors    int
		results   = make(chan error, totalOrders)
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	startTime := time.Now()

	// ======================
	// 3. Concurrent Requests
	// ======================
	for i := 0; i < totalOrders; i++ {
		wg.Add(1)
		go func(orderNum int) {
			defer wg.Done()

			_, err := orderService.CreateOrder(ctx, &pb.CreateOrderRequest{
				VariantsId: testVariantID,
				Quantity:   1,
				Total:      9.99,
				UserId:     fmt.Sprintf("user-%d", orderNum),
			})

			results <- err
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)

	// ======================
	// 4. Process Results
	// ======================
	for err := range results {
		if err != nil {
			errors++
			assert.Contains(t, err.Error(), "insufficient stock",
				"Unexpected error type")
		} else {
			successes++
		}
	}

	// ======================
	// 5. Assertions
	// ======================
	t.Logf("Test executed in %v", time.Since(startTime))
	t.Logf("Successes: %d, Errors: %d", successes, errors)

	// Verify stock be consistency
	var finalStock int
	err := db.Get(&finalStock,
		"SELECT quantity FROM stock WHERE variant_id = $1 FOR UPDATE",
		testVariantID,
	)
	require.NoError(t, err, "Failed to fetch final stock")

	assert.Equal(t, 0, finalStock, "Final stock should be zero")
	assert.Equal(t, initialStock, successes, "Number of successes should match initial stock")
	assert.Equal(t, totalOrders-initialStock, errors, "Number of errors should match overstock requests")
}
