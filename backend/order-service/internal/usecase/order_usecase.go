package usecase

import (
	"context"
	"time"

	"order-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type orderUseCase struct {
	orderRepo domain.OrderRepository
}

func NewOrderUseCase(orderRepo domain.OrderRepository) domain.OrderUseCase {
	return &orderUseCase{
		orderRepo: orderRepo,
	}
}

func (u *orderUseCase) CreateOrder(ctx context.Context, order *domain.Order) error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.Status = "pending"
	return u.orderRepo.Create(ctx, order)
}

func (u *orderUseCase) GetOrder(ctx context.Context, id primitive.ObjectID) (*domain.Order, error) {
	return u.orderRepo.GetByID(ctx, id)
}

func (u *orderUseCase) GetOrders(ctx context.Context, userID string) ([]domain.Order, error) {
	return u.orderRepo.GetAll(ctx, userID)
}

func (u *orderUseCase) UpdateOrder(ctx context.Context, order *domain.Order) error {
	order.UpdatedAt = time.Now()
	return u.orderRepo.Update(ctx, order)
}

func (u *orderUseCase) DeleteOrder(ctx context.Context, id primitive.ObjectID) error {
	return u.orderRepo.Delete(ctx, id)
}
