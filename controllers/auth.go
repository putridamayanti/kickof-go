package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kickof/models"
	"kickof/services"
	"net/http"
	"os"
	"time"
)

func SignUp(c *gin.Context) {
	var request models.Register

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	params := models.Register{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	url := os.Getenv("FRONTEND_URL") + "/activate/"
	token, err := services.Register(params, url)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	user := services.GetUser(bson.M{"email": request.Email}, nil)

	result := models.AuthResult{
		Token: *token,
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}

	c.JSON(http.StatusOK, models.Response{Data: result})
	return
}

func SignIn(c *gin.Context) {
	var request models.Login

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	token, err := services.SignIn(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	user := services.GetUser(bson.M{"email": request.Email}, nil)

	result := models.AuthResult{
		Token: *token,
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}

	c.JSON(http.StatusOK, models.Response{Data: result})
	return
}

func RefreshToken(c *gin.Context) {
	token, err := services.CheckHeader(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	email, err := services.VerifyToken(token)
	if err != nil {
		newToken, errGenerate := services.GenerateToken(email)
		if errGenerate != nil {
			c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
			return
		}

		token = *newToken
	}

	user := services.GetUser(bson.M{"email": email}, options.FindOne().SetProjection(bson.M{"password": 0}))
	if user == nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: token})
}

func Activate(c *gin.Context) {
	token := c.Param("token")

	_, err := services.Activate(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Activation Success"})
	return
}

func ForgotPassword(c *gin.Context) {
	var request models.ForgotPasswordRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	url := os.Getenv("FRONTEND_URL") + "/reset-password/"
	_, err = services.ForgotPassword(request.Email, url)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Reset email has been sent"})
	return
}

func UpdatePassword(c *gin.Context) {
	var request models.UpdatePasswordRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	_, err = services.UpdatePassword(request.Token, request.Password)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Success"})
	return
}

func GetProfile(c *gin.Context) {
	result := services.GetCurrentUser(c.Request)

	if result == nil {
		c.JSON(http.StatusNotFound, models.Response{Data: "User Not Found"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: result})
	return
}

func UpdateProfile(c *gin.Context) {
	res := services.GetCurrentUser(c.Request)
	if res == nil {
		c.JSON(http.StatusNotFound, models.Response{Data: "User Not Found"})
		return
	}

	var request models.User

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Data: err.Error()})
		return
	}

	request.Id = res.Id
	request.LastActive = time.Now()

	_, err = services.UpdateUser(res.Id, request)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Success"})
	return
}
