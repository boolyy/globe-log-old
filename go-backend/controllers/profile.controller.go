package controllers

import (
	"fmt"
	"net/http"

	"github.com/boolyy/globe-log/go-backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type friendReqInfo struct {
	Username   string `json:"username"`
	FriendName string `json:"friend"`
}

type privacyReqInfo struct {
	Username string `json:"username"`
	Privacy  string `json:"privacy"`
}

func validateFriendInfo(friendReqInfo friendReqInfo) error {
	if friendReqInfo.Username == "" {
		return fmt.Errorf("missing username")
	}

	if friendReqInfo.FriendName == "" {
		return fmt.Errorf("missing friend name")
	}

	return nil
}

func validatePrivacyInfo(privacyReqInfo privacyReqInfo) error {
	if privacyReqInfo.Username == "" {
		return fmt.Errorf("missing username")
	}

	if privacyReqInfo.Privacy == "" {
		return fmt.Errorf("missing privacy option")
	}

	return nil
}

//Append friend to friends array with set push, and also add friend to other person's friends array
func (profileController *Controller) AddFriend(context *gin.Context) {

	var friendReqInfo friendReqInfo

	if err := context.ShouldBindJSON(&friendReqInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := validateFriendInfo(friendReqInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	//Check if friend exists
	_, err := profileController.UserService.GetUser(friendReqInfo.FriendName)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "error finding user"})
		return
	}

	//Update user's friend list
	userUpdate := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "friends", Value: friendReqInfo.FriendName}}}}
	filter := bson.D{{Key: "username", Value: friendReqInfo.Username}}
	_, err = profileController.UserService.UpdateUser(filter, userUpdate)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	filter = bson.D{{Key: "username", Value: friendReqInfo.FriendName}}

	//Update friend's friend list with current user's username
	friendUpdate := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "friends", Value: friendReqInfo.Username}}}}
	_, err = profileController.UserService.UpdateUser(filter, friendUpdate)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusAccepted, gin.H{"message": "friend added successfully"})
}

func (profileController *Controller) DeleteFriend(context *gin.Context) {
	var friendReqInfo friendReqInfo

	if err := context.ShouldBindJSON(&friendReqInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := validateFriendInfo(friendReqInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userFilter := bson.D{{Key: "username", Value: friendReqInfo.Username}}
	userUpdate := bson.D{{Key: "$pull", Value: bson.D{{Key: "friends", Value: friendReqInfo.FriendName}}}}
	_, err := profileController.UserService.UpdateUser(userFilter, userUpdate)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	friendFilter := bson.D{{Key: "username", Value: friendReqInfo.FriendName}}
	friendUpdate := bson.D{{Key: "$pull", Value: bson.D{{Key: "friends", Value: friendReqInfo.Username}}}}
	_, err = profileController.UserService.UpdateUser(friendFilter, friendUpdate)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusAccepted, gin.H{"message": "friend deleted successfully"})
}

func (profileController *Controller) UpdatePrivacy(context *gin.Context) {
	var privacyReqInfo privacyReqInfo

	if err := context.ShouldBindJSON(&privacyReqInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := validatePrivacyInfo(privacyReqInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if _, exists := models.PrivacyMap[privacyReqInfo.Privacy]; !exists {
		context.JSON(http.StatusBadRequest, gin.H{"message": "incorrect privacy option"})
		return
	}

	filter := bson.D{{Key: "username", Value: privacyReqInfo.Username}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "privacy", Value: privacyReqInfo.Privacy}}}}
	updateResult, err := profileController.UserService.UpdateUser(filter, update)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, updateResult)
}

func (profileController *Controller) RegisterProfileRoutes(routerGroup *gin.RouterGroup) {
	profileRoute := routerGroup.Group("/profile")
	profileRoute.PUT("/friend", profileController.AddFriend)
	profileRoute.DELETE("/friend", profileController.DeleteFriend)
	profileRoute.PATCH("/privacy", profileController.UpdatePrivacy)
}
