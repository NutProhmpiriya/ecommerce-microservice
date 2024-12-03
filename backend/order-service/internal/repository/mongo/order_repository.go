package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"order-service/internal/domain"
)

type mongoOrderRepository struct {
	collection *mongo.Collection
}

func NewMongoOrderRepository(collection *mongo.Collection) domain.OrderRepository {
	return &mongoOrderRepository{
		collection: collection,
	}
}

func (r *mongoOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	_, err := r.collection.InsertOne(ctx, order)
	return err
}

func (r *mongoOrderRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Order, error) {
	var order domain.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *mongoOrderRepository) GetAll(ctx context.Context, userID string) ([]domain.Order, error) {
	filter := bson.M{}
	if userID != "" {
		filter["user_id"] = userID
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []domain.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *mongoOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	update := bson.M{
		"$set": bson.M{
			"product_id":  order.ProductID,
			"quantity":    order.Quantity,
			"total_price": order.TotalPrice,
			"status":      order.Status,
			"updated_at":  order.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": order.ID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *mongoOrderRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
