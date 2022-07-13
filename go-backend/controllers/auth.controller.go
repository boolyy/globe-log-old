package controllers

import (
	"net/http"

	"github.com/boolyy/globe-log/go-backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type signInInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (authController *Controller) RegisterUser(context *gin.Context) {

	var signInInfo signInInfo

	if err := context.ShouldBindJSON(&signInInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if signInInfo.Username == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "missing username"})
		return
	}

	if signInInfo.Password == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "missing password"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signInInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	signInInfo.Password = string(hashedPassword)

	//Create user based on inputted info
	user := models.User{
		Username:      signInInfo.Username,
		Password:      signInInfo.Password,
		PrivacyOption: models.PrivacyOption_Private,
		Friends:       []string{},
		Locations:     map[string]models.Location{},
	}

	if err = authController.UserService.CreateUser(user); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusAccepted, user)
}

func (authController *Controller) LoginUser(context *gin.Context) {
	//it's gonnna be a post. Look for username.
	//put the object in

	var loginInfo models.User
	var resultUser models.User

	if err := context.BindJSON(&loginInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if loginInfo.Username == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "missing username"})
		return
	}

	if loginInfo.Password == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "missing password"})
		return
	}

	//Look for username in mongo database
	resultUser, err := authController.UserService.GetUser(loginInfo.Username)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "user does not exist"})
		return
	}

	//Check if password matches
	hashedPassword := []byte(resultUser.Password)

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(loginInfo.Password)); err != nil {
		context.JSON(500, gin.H{"message": "wrong password"})
		return
	}

	context.JSON(http.StatusOK, resultUser)
}

func (authController *Controller) RegisterAuthRoutes(routerGroup *gin.RouterGroup) {
	authRoute := routerGroup.Group("/")
	authRoute.POST("/register", authController.RegisterUser)
	authRoute.POST("/login", authController.LoginUser)
}
