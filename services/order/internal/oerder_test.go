package handler

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wafi04/golang-backend/grpc/pb"
)

type mockDB struct {
	mock.Mock
}

func (m *mockDB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sqlx.Tx), args.Error(1)
}

func (m *mockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	mockArgs := m.Called(ctx, query, args)
	return mockArgs.Get(0).(*sql.Row)
}

func (m *mockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	mockArgs := m.Called(ctx, query, args)
	return mockArgs.Get(0).(sql.Result), mockArgs.Error(1)
}

type mockRedis struct {
	mock.Mock
}

func (m *mockRedis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *mockRedis) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *mockRedis) HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {
	args := m.Called(ctx, key, values)
	return args.Get(0).(*redis.BoolCmd)
}

type mockStockService struct {
	mock.Mock
}

func (m *mockStockService) GetStock(ctx context.Context, variantID string) (int64, error) {
	args := m.Called(ctx, variantID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockStockService) UpdateStock(ctx context.Context, variantID string, quantity int64) error {
	args := m.Called(ctx, variantID, quantity)
	return args.Error(0)
}

func TestCreateOrder_Success(t *testing.T) {
	mockDB := &sqlx.DB{}
	mockRedis := &redis.Client{}
	mockStock := new(mockStockService)

	mockStock.On("GetStock", mock.Anything, "test-variant-123").Return(int64(100), nil)
	mockStock.On("UpdateStock", mock.Anything, "test-variant-123", int64(5)).Return(nil)

	repo := NewOrderService(mockDB, mockRedis)

	req := &pb.CreateOrderRequest{
		VariantsId: "test-variant-123",
		Quantity:   5,
		Total:      99.99,
		UserId:     "user-123",
	}

	order, err := repo.CreateOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	mockStock.AssertExpectations(t)
}

func TestCreateOrder_InsufficientStock(t *testing.T) {
	mockDB := new(mockDB)
	mockRedis := new(mockRedis)

	orderService := &OrderRepository{
		db:          sqlx.NewDb(&sql.DB{}, "mock"),
		redisClient: &redis.Client{},
	}

	mockRedis.On("SetNX", mock.Anything, "order-lock:test-variant-456", "locked", lockTimeout).
		Return(redis.NewBoolResult(true, nil))

	mockDB.On("BeginTxx", mock.Anything, mock.Anything).
		Return(&sqlx.Tx{}, nil)

	mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).
		Return(&sql.Row{}).Run(func(args mock.Arguments) {
	})

	req := &pb.CreateOrderRequest{
		VariantsId: "test-variant-456",
		Quantity:   5,
		Total:      100.50,
		UserId:     "test-user-789",
	}

	order, err := orderService.CreateOrder(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "insufficient stock")

	mockRedis.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestCreateOrder_LockAcquisitionFailed(t *testing.T) {
	mockDB := new(mockDB)
	mockRedis := new(mockRedis)

	orderService := &OrderRepository{
		db:          sqlx.NewDb(&sql.DB{}, "mock"),
		redisClient: &redis.Client{},
	}

	mockRedis.On("SetNX", mock.Anything, "order-lock:test-variant-456", "locked", lockTimeout).
		Return(redis.NewBoolResult(false, errors.New("lock acquisition failed")))

	req := &pb.CreateOrderRequest{
		VariantsId: "test-variant-456",
		Quantity:   5,
		Total:      100.50,
		UserId:     "test-user-789",
	}

	order, err := orderService.CreateOrder(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "failed to acquire lock")

	mockRedis.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}
