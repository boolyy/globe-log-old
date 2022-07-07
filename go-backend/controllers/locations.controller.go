package controllers

import (
	"fmt"
	"net/http"

	"github.com/boolyy/globe-log/go-backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Key of location in user's locations map in mongodb. Format of "(lat, long)"
type DeleteLocationReqInfo struct {
	Username string `json:"username"`
	Key      string `json:"locationKey"`
}

type LocationReqInfo struct {
	Username string          `json:"username"`
	Location models.Location `json:"location"`
}

func createLocationKeyFromCords(cords []float32) string {
	return "(" + fmt.Sprint(cords[0]) + "," + fmt.Sprint(cords[1]) + ")"
}

func areCordsValid(cords []float32) error {

	latitude, longitude := cords[0], cords[1]

	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("latitude not in range: -90 to 90")
	}

	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("longitude not in range: -180 to 180")
	}

	return nil
}

func (locationController *Controller) AddLocation(context *gin.Context) {
	var locationReqInfo LocationReqInfo

	if err := context.ShouldBindJSON(&locationReqInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	//validation
	if err := areCordsValid(locationReqInfo.Location.Coordinates); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	//create filter update
	locationKey := createLocationKeyFromCords(locationReqInfo.Location.Coordinates)
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "locations." + locationKey, Value: locationReqInfo.Location}}}}

	//Update user of locationReqInfo.Username
	updateResult, err := locationController.UserService.UpdateUser(locationReqInfo.Username, update)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, updateResult)
}

func (locationController *Controller) UpdateLocation(context *gin.Context) {

}

func (locationController *Controller) DeleteLocation(context *gin.Context) {
	var deleteLocationReqInfo DeleteLocationReqInfo

	if err := context.ShouldBindJSON(&deleteLocationReqInfo); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	username, locationToDelete := deleteLocationReqInfo.Username, deleteLocationReqInfo.Key
	// call update uset with unset or something to delete location key for given location
	update := bson.D{primitive.E{Key: "$unset", Value: bson.D{primitive.E{Key: "locations." + locationToDelete, Value: ""}}}}

	updateResult, err := locationController.UserService.UpdateUser(username, update)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, updateResult)
}

func (locationController *Controller) RegisterLocationRoutes(routerGroup *gin.RouterGroup) {
	locationRoute := routerGroup.Group("/location")
	locationRoute.PUT("", locationController.AddLocation)
	locationRoute.DELETE("", locationController.DeleteLocation)
}
