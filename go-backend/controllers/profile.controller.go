package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type friendReqInfo struct {
	Username   string `json:"username"`
	FriendName string `json:"friend"`
}

//Append friend to friends array with set push, and also add friend to other person's friends array
func (profileController *Controller) AddFriend(context *gin.Context) {

	var friendReqInfo friendReqInfo

	if err := context.ShouldBindJSON(&friendReqInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if friendReqInfo.Username == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "missing username"})
		return
	}

	if friendReqInfo.FriendName == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "missing friend to be added"})
		return
	}

	//Check if friend exists
	_, err := profileController.UserService.GetUser(friendReqInfo.FriendName)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "error finding user"})
		return
	}

	//Update user's friend list
	userUpdate := bson.D{primitive.E{Key: "$addToSet", Value: bson.D{primitive.E{Key: "friends", Value: friendReqInfo.FriendName}}}}
	_, err = profileController.UserService.UpdateUser(friendReqInfo.Username, userUpdate)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	//Update friend's friend list with current user's username
	friendUpdate := bson.D{primitive.E{Key: "$addToSet", Value: bson.D{primitive.E{Key: "friends", Value: friendReqInfo.Username}}}}
	_, err = profileController.UserService.UpdateUser(friendReqInfo.FriendName, friendUpdate)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusAccepted, gin.H{"message": "friend added successfully"})
}

func (profileController *Controller) DeleteFriend(context *gin.Context) {

}

func (profileController *Controller) UpdatePrivacy(context *gin.Context) {

}

func (profileController *Controller) RegisterProfileRoutes(routerGroup *gin.RouterGroup) {
	profileRoute := routerGroup.Group("/profile")
	profileRoute.PUT("", profileController.AddFriend)
}
