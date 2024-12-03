package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"kickof/models"
	"kickof/services"
	"net/http"
	"time"
)

func GetProjects(c *gin.Context) {
	var query models.Query

	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	filters := query.GetQueryFind()

	if query.UserId != "" {
		filters["userIds"] = bson.M{
			"$in": []string{query.UserId},
		}
	}

	opts := query.GetOptions()

	results := services.GetProjectsWithPagination(filters, opts, query)

	c.JSON(http.StatusOK, models.Response{Data: results})
	return
}

func CreateProject(c *gin.Context) {
	var request models.Project

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, models.Response{Data: err.Error()})
		return
	}

	request.Id = uuid.New().String()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	_, err = services.CreateProject(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: request})
	return
}

func GetProjectByIdOrCode(c *gin.Context) {
	id := c.Param("id")

	result := services.GetProject(bson.M{"id": id}, nil)
	if result == nil {
		result = services.GetProject(bson.M{"code": id}, nil)
		if result == nil {
			c.JSON(http.StatusNotFound, models.Result{Data: "Data Not Found"})
			return
		}
	}

	c.JSON(http.StatusOK, models.Response{Data: result})
}

func UpdateProject(c *gin.Context) {
	id := c.Param("id")

	data := services.GetProject(bson.M{"id": id}, nil)
	if data == nil {
		c.JSON(http.StatusNotFound, models.Result{Data: "Data Not Found"})
		return
	}

	var request models.Project

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	_, err = services.UpdateProject(id, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(200, models.Response{Data: request})
}

func DeleteProject(c *gin.Context) {
	id := c.Param("id")

	_, err := services.DeleteProject(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: "Failed Delete Data"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Success"})
}
