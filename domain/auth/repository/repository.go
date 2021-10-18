package repository

import (
	"barbar/domain/auth/entity"
	"barbar/pkg/redis"
	"barbar/pkg/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthRepositoryInterface interface {
	Create(ctx context.Context, user entity.Auth) (entity.Auth, error)
	CreateToRedis(ctx context.Context, auth entity.Auth) (err error)
	Update(ctx context.Context, user entity.Auth) (entity.Auth, error)
	DeleteDataRedis(ctx context.Context, userId string) error

	GetById(ctx context.Context, id string) (*entity.Auth, error)
	GetAllAuth(ctx context.Context, search string, sortBy map[string]bool, limit, offset int64) ([]entity.Auth, error)
}

type userRepository struct {
	db    *mongo.Database
	cache *redis.Redis
}

func (u userRepository) Create(ctx context.Context, auth entity.Auth) (newAuth entity.Auth, err error) {
	_, err = u.db.Collection("auth").InsertOne(ctx, auth)
	return auth, err
}

func (u userRepository) CreateToRedis(ctx context.Context, auth entity.Auth) (err error) {
	err = u.cache.SetDataRedis(ctx, "users-"+auth.UserID, auth.AccessToken)
	return err
}

func (u userRepository) Update(ctx context.Context, auth entity.Auth) (newAuth entity.Auth, err error) {
	_, err = u.db.Collection("users").UpdateOne(ctx, bson.M{"uid": auth.UID}, bson.D{
		{"$set", auth},
	},
	)
	if err != nil {
		return auth, err
	}

	return auth, nil
}

func (u userRepository) DeleteDataRedis(ctx context.Context, userId string) error {
	err := u.cache.DeleteDataRedis(ctx, "users-"+userId)
	return err
}

func (u userRepository) GetById(ctx context.Context, id string) (*entity.Auth, error) {
	panic("implement me")
}

func (u userRepository) GetAllAuth(ctx context.Context, search string,
	sortBy map[string]bool, limit, offset int64) (auths []entity.Auth, err error) {
	findOptions := options.Find()
	findOptions.SetSort(utils.MongoSetSortFromMap(sortBy))
	findOptions.SetLimit(limit)
	findOptions.SetSkip((offset) * limit)

	cursor, err := u.db.Collection("auth").Find(ctx, bson.M{"userid": search}, findOptions)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &auths)
	return auths, err
}

func NewAuthRepository(db *mongo.Database, cache *redis.Redis) *userRepository {
	return &userRepository{
		db:    db,
		cache: cache,
	}
}
