package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"order-service/internal/domain"
	mockRepo "order-service/internal/repository/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateOrder(t *testing.T) {
	mockRepo := new(mockRepo.MockOrderRepository)
	useCase := NewOrderUseCase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		order := &domain.Order{
			UserID:     "123",
			ProductID:  "456",
			Quantity:   2,
			TotalPrice: 1000,
		}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(o *domain.Order) bool {
			return o.UserID == order.UserID &&
				o.ProductID == order.ProductID &&
				o.Quantity == order.Quantity &&
				o.TotalPrice == order.TotalPrice &&
				o.Status == "pending" &&
				!o.CreatedAt.IsZero() &&
				!o.UpdatedAt.IsZero()
		})).Return(nil).Once()

		err := useCase.CreateOrder(context.Background(), order)

		assert.NoError(t, err)
		assert.Equal(t, "pending", order.Status)
		assert.NotZero(t, order.CreatedAt)
		assert.NotZero(t, order.UpdatedAt)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		order := &domain.Order{
			UserID:     "123",
			ProductID:  "456",
			Quantity:   2,
			TotalPrice: 1000,
		}

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError).Once()

		err := useCase.CreateOrder(context.Background(), order)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetOrder(t *testing.T) {
	mockRepo := new(mockRepo.MockOrderRepository)
	useCase := NewOrderUseCase(mockRepo)

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

		mockRepo.On("GetByID", mock.Anything, id).Return(expectedOrder, nil).Once()

		order, err := useCase.GetOrder(context.Background(), id)

		assert.NoError(t, err)
		assert.Equal(t, expectedOrder, order)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		id := primitive.NewObjectID()
		mockRepo.On("GetByID", mock.Anything, id).Return(nil, nil).Once()

		order, err := useCase.GetOrder(context.Background(), id)

		assert.NoError(t, err)
		assert.Nil(t, order)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		id := primitive.NewObjectID()
		mockRepo.On("GetByID", mock.Anything, id).Return(nil, assert.AnError).Once()

		order, err := useCase.GetOrder(context.Background(), id)

		assert.Error(t, err)
		assert.Nil(t, order)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetOrders(t *testing.T) {
	mockRepo := new(mockRepo.MockOrderRepository)
	useCase := NewOrderUseCase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		userID := "123"
		expectedOrders := []domain.Order{
			{
				ID:         primitive.NewObjectID(),
				UserID:     userID,
				ProductID:  "456",
				Quantity:   2,
				TotalPrice: 1000,
				Status:     "pending",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				ID:         primitive.NewObjectID(),
				UserID:     userID,
				ProductID:  "789",
				Quantity:   1,
				TotalPrice: 500,
				Status:     "completed",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		}

		mockRepo.On("GetAll", mock.Anything, userID).Return(expectedOrders, nil).Once()

		orders, err := useCase.GetOrders(context.Background(), userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedOrders, orders)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		userID := "123"
		mockRepo.On("GetAll", mock.Anything, userID).Return(nil, assert.AnError).Once()

		orders, err := useCase.GetOrders(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, orders)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateOrder(t *testing.T) {
	mockRepo := new(mockRepo.MockOrderRepository)
	useCase := NewOrderUseCase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		order := &domain.Order{
			ID:         primitive.NewObjectID(),
			UserID:     "123",
			ProductID:  "456",
			Quantity:   3,
			TotalPrice: 1500,
			Status:     "processing",
		}

		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(o *domain.Order) bool {
			return o.ID == order.ID &&
				o.UserID == order.UserID &&
				o.ProductID == order.ProductID &&
				o.Quantity == order.Quantity &&
				o.TotalPrice == order.TotalPrice &&
				o.Status == order.Status &&
				!o.UpdatedAt.IsZero()
		})).Return(nil).Once()

		err := useCase.UpdateOrder(context.Background(), order)

		assert.NoError(t, err)
		assert.NotZero(t, order.UpdatedAt)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		order := &domain.Order{
			ID:         primitive.NewObjectID(),
			UserID:     "123",
			ProductID:  "456",
			Quantity:   3,
			TotalPrice: 1500,
			Status:     "processing",
		}

		mockRepo.On("Update", mock.Anything, mock.Anything).Return(assert.AnError).Once()

		err := useCase.UpdateOrder(context.Background(), order)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteOrder(t *testing.T) {
	mockRepo := new(mockRepo.MockOrderRepository)
	useCase := NewOrderUseCase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		id := primitive.NewObjectID()
		mockRepo.On("Delete", mock.Anything, id).Return(nil).Once()

		err := useCase.DeleteOrder(context.Background(), id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		id := primitive.NewObjectID()
		mockRepo.On("Delete", mock.Anything, id).Return(assert.AnError).Once()

		err := useCase.DeleteOrder(context.Background(), id)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
