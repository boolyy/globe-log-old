package services

import (
	"github.com/boolyy/globe-log/go-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//All the methods to interact with mongodb
type UserService interface {
	CreateUser(models.User) error
	GetUser(string) (models.User, error)
	GetField(bson.D) *mongo.SingleResult
	GetAll() ([]*models.User, error)
	UpdateUser(bson.D, bson.D) (*mongo.UpdateResult, error)
	DeleteUser(string) error
}
