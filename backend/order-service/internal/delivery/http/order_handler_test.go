package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"order-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockOrderUseCase struct {
	mock.Mock
}

func (m *MockOrderUseCase) CreateOrder(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderUseCase) GetOrder(ctx context.Context, id primitive.ObjectID) (*domain.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderUseCase) GetOrders(ctx context.Context, userID string) ([]domain.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Order), args.Error(1)
}

func (m *MockOrderUseCase) UpdateOrder(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderUseCase) DeleteOrder(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateOrder(t *testing.T) {
	mockUseCase := new(MockOrderUseCase)
	router := mux.NewRouter()
	NewOrderHandler(router, mockUseCase)

	t.Run("Success", func(t *testing.T) {
		order := domain.Order{
			UserID:     "123",
			ProductID:  "456",
			Quantity:   2,
			TotalPrice: 1000,
		}

		mockUseCase.On("CreateOrder", mock.Anything, mock.MatchedBy(func(o *domain.Order) bool {
			return o.UserID == order.UserID &&
				o.ProductID == order.ProductID &&
				o.Quantity == order.Quantity &&
				o.TotalPrice == order.TotalPrice
		})).Return(nil).Once()

		body, _ := json.Marshal(order)
		req := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("UseCase Error", func(t *testing.T) {
		order := domain.Order{
			UserID:     "123",
			ProductID:  "456",
			Quantity:   2,
			TotalPrice: 1000,
		}

		mockUseCase.On("CreateOrder", mock.Anything, mock.Anything).Return(assert.AnError).Once()

		body, _ := json.Marshal(order)
		req := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestGetOrder(t *testing.T) {
	mockUseCase := new(MockOrderUseCase)
	router := mux.NewRouter()
	NewOrderHandler(router, mockUseCase)

	t.Run("Success", func(t *testing.T) {
		id := primitive.NewObjectID()
		expectedOrder := &domain.Order{
			ID:         id,
			UserID:     "123",
			ProductID:  "456",
			Quantity:   2,
			TotalPrice: 1000,
			Status:     "pending",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		mockUseCase.On("GetOrder", mock.Anything, id).Return(expectedOrder, nil).Once()

		req := httptest.NewRequest("GET", "/api/v1/orders/"+id.Hex(), nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response domain.Order
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedOrder.ID, response.ID)
		assert.Equal(t, expectedOrder.UserID, response.UserID)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/orders/invalid-id", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		id := primitive.NewObjectID()
		mockUseCase.On("GetOrder", mock.Anything, id).Return(nil, nil).Once()

		req := httptest.NewRequest("GET", "/api/v1/orders/"+id.Hex(), nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestGetOrders(t *testing.T) {
	mockUseCase := new(MockOrderUseCase)
	router := mux.NewRouter()
	NewOrderHandler(router, mockUseCase)

	t.Run("Success", func(t *testing.T) {
		userID := "123"
		expectedOrders := []domain.Order{
			{
				ID:         primitive.NewObjectID(),
				UserID:     userID,
				ProductID:  "456",
				Quantity:   2,
				TotalPrice: 1000,
			},
		}

		mockUseCase.On("GetOrders", mock.Anything, userID).Return(expectedOrders, nil).Once()

		req := httptest.NewRequest("GET", "/api/v1/orders?user_id="+userID, nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response []domain.Order
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedOrders), len(response))
		mockUseCase.AssertExpectations(t)
	})
}
