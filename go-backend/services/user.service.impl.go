package services

import (
	"context"

	"github.com/boolyy/globe-log/go-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	userCollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(userCollection *mongo.Collection, mongoContext context.Context) UserService {
	return &UserServiceImpl{
		userCollection: userCollection,
		ctx:            mongoContext,
	}
}

func (userServiceImpl *UserServiceImpl) CreateUser(user models.User) error {
	_, err := userServiceImpl.userCollection.InsertOne(userServiceImpl.ctx, user)
	return err
}

func (userServiceImpl *UserServiceImpl) GetUser(username string) (models.User, error) {
	var resultUser models.User
	query := bson.D{{Key: "username", Value: username}}
	err := userServiceImpl.userCollection.FindOne(userServiceImpl.ctx, query).Decode(&resultUser)
	return resultUser, err
}

func (userServiceImpl *UserServiceImpl) UpdateUser(username string, update bson.D) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "username", Value: username}}
	updateResult, err := userServiceImpl.userCollection.UpdateOne(userServiceImpl.ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updateResult, nil
}

func (userServiceImpl *UserServiceImpl) DeleteUser(username string) error {
	filter := bson.D{{Key: "username", Value: username}}
	_, err := userServiceImpl.userCollection.DeleteOne(userServiceImpl.ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (UserServiceImpl *UserServiceImpl) GetAll() ([]*models.User, error) {
	return nil, nil
}
