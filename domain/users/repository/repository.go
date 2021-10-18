package repository

import (
	"barbar/domain/users/entity"
	"barbar/pkg/redis"
	"barbar/pkg/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, userId string) error

	GetById(ctx context.Context, id string) (*entity.User, error)
	GetAllUser(ctx context.Context, search string, sortBy map[string]bool, limit, offset int64) ([]entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type userRepository struct {
	db    *mongo.Database
	cache *redis.Redis
}

func (u userRepository) Create(ctx context.Context, user entity.User) (newUser entity.User, err error) {
	_, err = u.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return newUser, err
	}

	return user, nil
}

func (u userRepository) Update(ctx context.Context, user entity.User) (newUser entity.User, err error) {
	_, err = u.db.Collection("users").UpdateOne(ctx, bson.M{"uid": user.UID}, bson.D{
		{"$set", user},
	},
	)
	if err != nil {
		return newUser, err
	}

	return user, nil
}

func (u userRepository) Delete(ctx context.Context, userId string) (err error) {
	_, err = u.db.Collection("users").DeleteOne(ctx, bson.M{"uid": userId})
	return err
}

func (u userRepository) GetById(ctx context.Context, id string) (user *entity.User, err error) {
	err = u.db.Collection("users").FindOne(ctx, bson.M{"uid": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return user, err
	}

	return user, nil
}

func (u userRepository) GetAllUser(ctx context.Context,
	search string, sortBy map[string]bool, limit, offset int64) (users []entity.User, err error) {
	findOptions := options.Find()
	findOptions.SetSort(utils.MongoSetSortFromMap(sortBy))
	findOptions.SetLimit(limit)
	findOptions.SetSkip((offset) * limit)

	cursor, err := u.db.Collection("users").Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &users)
	return users, err
}

func (u userRepository) GetByEmail(ctx context.Context, email string) (user *entity.User, err error) {
	err = u.db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return user, err
	}

	return user, nil
}

func NewUserRepository(db *mongo.Database, cache *redis.Redis) *userRepository {
	return &userRepository{
		db:    db,
		cache: cache,
	}
}
