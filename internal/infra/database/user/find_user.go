package user

import (
	"context"
	"errors"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/logger"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/entity/user_entity"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/internal_error"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type UserEntityMongo struct {
	Id   primitive.ObjectID `bson:"_id,omitempty" copier:"Id"`
	Name string             `bson:"name" copier:"Name"`
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: database.Collection("users"),
	}
}

func (repo *UserRepository) FindUserByID(ctx context.Context, userId string) (*user_entity.User, *internal_error.InternalError) {
	objectUserId, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": objectUserId}

	var userEntityMongo UserEntityMongo
	err := repo.Collection.FindOne(ctx, filter).Decode(&userEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error("User not found", err, zap.String("UserID", userId))
			return nil, internal_error.NewNotFoundError("User not found")
		}

		logger.Error("Error finding user_usecase by ID", err, zap.String("UserID", userId))
		return nil, internal_error.NewInternalServerError("Error finding user_usecase by ID")
	}

	var userEntity user_entity.User
	if errCopier := copier.Copy(&userEntity, &userEntityMongo); errCopier != nil {
		logger.Error("Error finding user_usecase by ID", errCopier, zap.String("UserID", userId))
		return nil, internal_error.NewInternalServerError("Error finding user_usecase by ID")
	}

	userEntity.Id = userId

	return &userEntity, nil
}
