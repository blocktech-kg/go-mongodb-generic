package mongodb

import (
	"context"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"
)

type CRUDDBService[T any] interface {
	// Create item in DB
	// Note: item ID used in database SHOULD BE set externally
	// if some failed, return err
	Create(ctx context.Context, item *T) (err error)

	// Get an item by id
	// if some failed, return err
	Get(ctx context.Context, id any) (item *T, err error)

	// Update an item identified by id
	// if some failed, return err
	Update(ctx context.Context, id any, item *T) (err error)

	// UpdateAttributes updates item's attributes 'attrs' identified by filter 'sels'
	// if some failed, return err
	UpdateAttributes(ctx context.Context, sels map[string]any, attrs map[string]any) error

	// Delete item in DB and identified by id
	// if some failed, return err
	Delete(ctx context.Context, id any) error

	// DeleteRange delete items in DB and identified by sels
	// if some failed, return err
	DeleteRange(ctx context.Context, sels map[string]any) error

	// ListAll uses for getting all items in DB for entity
	// if some failed, return err
	ListAll(ctx context.Context) (items []T, err error)

	// Find exact one item by sels filter (logical AND)
	// if some failed, return err
	Find(ctx context.Context, sels map[string]any) (item *T, err error)

	// Exists uses for checking if item exists with sels filter (logical AND)
	// if found return item and exist=true
	// if not found return nil and exist=false
	// if some failed, return err
	Exists(ctx context.Context, sels map[string]any) (item *T, exist bool, err error)

	// List all items by sels filter (logical AND)
	// if some failed, return err
	List(ctx context.Context, sels map[string]any) ([]T, error)

	// CreateIndex create index based on sels and unique flag
	// if some failed, return err
	CreateIndex(ctx context.Context, sels map[string]int, unique bool) (string, error)
}

func NewGenericObjectDBCtrl[T any](dbCollection *mongo.Collection) *genericObjectDBCtrl[T] {
	return &genericObjectDBCtrl[T]{
		db: dbCollection,
	}
}

type genericObjectDBCtrl[T any] struct {
	db *mongo.Collection
}

func (c *genericObjectDBCtrl[T]) Create(ctx context.Context, item *T) error {
	log.Debug("DB DEBUG: Started c.db.InsertOne(ctx, &item)")
	defer log.Debug("DB DEBUG: finished c.db.InsertOne(ctx, &item)")
	now := time.Now()
	createdAtField := reflect.ValueOf(item).Elem().FieldByName("CreatedAt")
	if createdAtField.IsValid() && createdAtField.CanSet() {
		createdAtField.Set(reflect.ValueOf(now))
	}
	updatedAtField := reflect.ValueOf(item).Elem().FieldByName("UpdatedAt")
	if updatedAtField.IsValid() && updatedAtField.CanSet() {
		updatedAtField.Set(reflect.ValueOf(now))
	}

	_, err := c.db.InsertOne(ctx, &item)
	if err != nil {
		return err
	}

	return nil
}

func (c *genericObjectDBCtrl[T]) Get(ctx context.Context, id any) (*T, error) {
	log.Debug("DB DEBUG: Started c.db.FindOne(ctx, filter)")
	defer log.Debug("DB DEBUG: finished c.db.FindOne(ctx, filter)")
	result := new(T)
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	err := c.db.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *genericObjectDBCtrl[T]) Find(ctx context.Context, sels map[string]any) (*T, error) {
	log.Debug("DB DEBUG: Started c.db.FindOne(ctx, filter)")
	defer log.Debug("DB DEBUG: finished c.db.FindOne(ctx, filter)")

	result := new(T)

	var filter bson.D
	for k, v := range sels {
		filter = append(filter, bson.E{k, v})
	}

	err := c.db.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *genericObjectDBCtrl[T]) Exists(ctx context.Context, sels map[string]any) (*T, bool, error) {
	log.Debug("DB DEBUG: Started c.Find(ctx, sels)")
	defer log.Debug("DB DEBUG: finished c.Find(ctx, sels)")

	result, err := c.Find(ctx, sels)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return result, true, nil
}

func (c *genericObjectDBCtrl[T]) Update(ctx context.Context, id any, item *T) error {
	log.Debug("DB DEBUG: Started c.db.UpdateOne")
	defer log.Debug("DB DEBUG: finished c.db.UpdateOne")
	now := time.Now()
	updatedAtField := reflect.ValueOf(item).Elem().FieldByName("UpdatedAt")
	if updatedAtField.IsValid() && updatedAtField.CanSet() {
		updatedAtField.Set(reflect.ValueOf(now))
	}
	dataByte, err := bson.Marshal(item)
	if err != nil {
		return err
	}

	var update bson.M
	err = bson.Unmarshal(dataByte, &update)
	if err != nil {
		return err
	}

	_, err = c.db.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			bson.E{Key: "$set", Value: update},
		},
	)

	if err != nil {
		return err
	}
	return nil
}

func (c *genericObjectDBCtrl[T]) UpdateAttributes(ctx context.Context, sels map[string]any, attrs map[string]any) error {
	log.Debug("DB DEBUG: Started c.db.UpdateMany")
	defer log.Debug("DB DEBUG: finished c.db.UpdateMany")
	var filter bson.D
	for k, v := range sels {
		filter = append(filter, bson.E{k, v})
	}

	var update bson.M
	attrs["updated_at"] = time.Now()
	dataByte, err := bson.Marshal(attrs)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(dataByte, &update)
	if err != nil {
		return err
	}

	_, err = c.db.UpdateMany(
		ctx,
		filter,
		bson.D{
			bson.E{Key: "$set", Value: update},
		},
	)

	if err != nil {
		return err
	}
	return nil
}

func (c *genericObjectDBCtrl[T]) Delete(ctx context.Context, id any) error {
	filter := bson.D{
		bson.E{Key: "_id", Value: id},
	}
	_, err := c.db.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (c *genericObjectDBCtrl[T]) DeleteRange(ctx context.Context, sels map[string]any) error {
	var filter bson.D
	for k, v := range sels {
		filter = append(filter, bson.E{k, v})
	}
	_, err := c.db.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (c *genericObjectDBCtrl[T]) ListAll(ctx context.Context) ([]T, error) {
	log.Debug("DB DEBUG: Started c.db.Find(ctx, filter)")
	defer log.Debug("DB DEBUG: finished c.db.Find(ctx, filter)")

	filter := bson.D{bson.E{}}

	cursor, err := c.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	results := []T{}
	for cursor.Next(ctx) {
		var result T
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}
	return results, nil
}

func (c *genericObjectDBCtrl[T]) List(ctx context.Context, sels map[string]any) ([]T, error) {
	log.Debug("DB DEBUG: Started c.db.Find(ctx, filter)")
	defer log.Debug("DB DEBUG: finished c.db.Find(ctx, filter)")
	var filter bson.D
	for k, v := range sels {
		filter = append(filter, bson.E{k, v})
	}

	cursor, err := c.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	results := []T{}
	for cursor.Next(ctx) {
		var result T
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}
	return results, nil
}

func (c *genericObjectDBCtrl[T]) CreateIndex(ctx context.Context, sels map[string]int, unique bool) (string, error) {
	var indexKeys bson.D
	for key, value := range sels {
		indexKeys = append(indexKeys, bson.E{Key: key, Value: value})
	}

	indexModel := mongo.IndexModel{
		Keys:    indexKeys,
		Options: options.Index().SetUnique(unique),
	}

	// Создание индекса
	indexName, err := c.db.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return "", err
	}

	return indexName, nil
}
