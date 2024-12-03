package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      string            `json:"user_id" bson:"user_id"`
	ProductID   string            `json:"product_id" bson:"product_id"`
	Quantity    int               `json:"quantity" bson:"quantity"`
	TotalPrice  float64           `json:"total_price" bson:"total_price"`
	Status      string            `json:"status" bson:"status"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Order, error)
	GetAll(ctx context.Context, userID string) ([]Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type OrderUseCase interface {
	CreateOrder(ctx context.Context, order *Order) error
	GetOrder(ctx context.Context, id primitive.ObjectID) (*Order, error)
	GetOrders(ctx context.Context, userID string) ([]Order, error)
	UpdateOrder(ctx context.Context, order *Order) error
	DeleteOrder(ctx context.Context, id primitive.ObjectID) error
}
