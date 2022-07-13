package controllers

import (
	"fmt"
	"net/http"

	"github.com/boolyy/globe-log/go-backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type deleteTripReqInfo struct {
	Username string `json:"username"`
	Key      string `json:"tripKey"`
}

type tripReqInfo struct {
	Username string      `json:"username"`
	Trip     models.Trip `json:"trip"`
}

func createTripKeyFromCords(tripReqInfo tripReqInfo) string {
	startCords, endCords := tripReqInfo.Trip.StartCoordinates, tripReqInfo.Trip.EndCoordinates
	return createLocationKeyFromCords(startCords) + "-" + createLocationKeyFromCords(endCords)
}

func (tripController *Controller) AddTrip(context *gin.Context) {
	var tripReqInfo tripReqInfo

	if err := context.ShouldBindJSON(&tripReqInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := tripController.validateTripInfo(tripReqInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tripKey := createTripKeyFromCords(tripReqInfo)
	startLocationKey := createLocationKeyFromCords(tripReqInfo.Trip.StartCoordinates)
	endLocationKey := createLocationKeyFromCords(tripReqInfo.Trip.EndCoordinates)

	tripDoesNotExistFilter := primitive.E{Key: "trips." + tripKey, Value: bson.D{{Key: "$exists", Value: false}}}
	startLocationExistsFilter := primitive.E{Key: "locations." + startLocationKey, Value: bson.D{{Key: "$exists", Value: true}}}
	endLocationExistsFilter := primitive.E{Key: "locations." + endLocationKey, Value: bson.D{{Key: "$exists", Value: true}}}

	filter := bson.D{{Key: "username", Value: tripReqInfo.Username}, tripDoesNotExistFilter, startLocationExistsFilter, endLocationExistsFilter}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "trips." + tripKey, Value: tripReqInfo.Trip}}}}

	updateResult, err := tripController.UserService.UpdateUser(filter, update)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, updateResult)
}

func (tripController *Controller) DeleteTrip(context *gin.Context) {
	var deleteTripReqInfo deleteTripReqInfo

	if err := context.ShouldBindJSON(&deleteTripReqInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	username, tripToDelete := deleteTripReqInfo.Username, deleteTripReqInfo.Key

	filter := bson.D{{Key: "username", Value: username}}
	update := bson.D{{Key: "$unset", Value: bson.D{{Key: "trips." + tripToDelete, Value: ""}}}}

	updateResult, err := tripController.UserService.UpdateUser(filter, update)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, updateResult)
}

func (tripController *Controller) UpdateTrip(context *gin.Context) {
	var tripReqInfo tripReqInfo

	if err := context.ShouldBindJSON(&tripReqInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := tripController.validateTripInfo(tripReqInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tripKey := createTripKeyFromCords(tripReqInfo)
	startLocationKey := createLocationKeyFromCords(tripReqInfo.Trip.StartCoordinates)
	endLocationKey := createLocationKeyFromCords(tripReqInfo.Trip.EndCoordinates)

	tripDoesNotExistFilter := primitive.E{Key: "trips." + tripKey, Value: bson.D{{Key: "$exists", Value: true}}}
	startLocationExistsFilter := primitive.E{Key: "locations." + startLocationKey, Value: bson.D{{Key: "$exists", Value: true}}}
	endLocationExistsFilter := primitive.E{Key: "locations." + endLocationKey, Value: bson.D{{Key: "$exists", Value: true}}}

	filter := bson.D{{Key: "username", Value: tripReqInfo.Username}, tripDoesNotExistFilter, startLocationExistsFilter, endLocationExistsFilter}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "trips." + tripKey, Value: tripReqInfo.Trip}}}}

	updateResult, err := tripController.UserService.UpdateUser(filter, update)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, updateResult)
}

func (tripController *Controller) validateTripLocations(tripReqInfo tripReqInfo) error {
	// validate first location
	startLocationKey := createLocationKeyFromCords(tripReqInfo.Trip.StartCoordinates)

	startLocationFilter := bson.D{{Key: "username", Value: tripReqInfo.Username}, {Key: "locations." + startLocationKey, Value: bson.D{{Key: "$exists", Value: true}}}}
	if err := tripController.UserService.GetField(startLocationFilter).Err(); err != nil {
		return fmt.Errorf("error finding start location:" + err.Error())
	}

	endLocationKey := createLocationKeyFromCords(tripReqInfo.Trip.EndCoordinates)

	endLocationFilter := bson.D{{Key: "username", Value: tripReqInfo.Username}, {Key: "locations." + endLocationKey, Value: bson.D{{Key: "$exists", Value: true}}}}
	if err := tripController.UserService.GetField(endLocationFilter).Err(); err != nil {
		return fmt.Errorf("error finding end location")
	}

	return nil
}

func (tripController *Controller) validateTripInfo(tripReqInfo tripReqInfo) error {
	if tripReqInfo.Username == "" {
		return fmt.Errorf("missing username")
	}

	if err := areCordsValid(tripReqInfo.Trip.StartCoordinates); err != nil {
		return fmt.Errorf("error with start location: " + err.Error())
	}

	if err := areCordsValid(tripReqInfo.Trip.EndCoordinates); err != nil {
		return fmt.Errorf("error with end location: " + err.Error())
	}

	if tripReqInfo.Trip.Title == "" {
		return fmt.Errorf("missing title")
	}

	// Validate that cords given for trip already exists
	if err := tripController.validateTripLocations(tripReqInfo); err != nil {
		return err
	}

	return nil
}

func (tripController *Controller) RegisterTripRoutes(routerGroup *gin.RouterGroup) {
	tripRoute := routerGroup.Group("/trip")
	tripRoute.PUT("", tripController.AddTrip)
	tripRoute.DELETE("", tripController.DeleteTrip)
	tripRoute.PATCH("", tripController.UpdateTrip)
}
