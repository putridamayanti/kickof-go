package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kickof/models"
	"kickof/services"
	"net/http"
	"time"
)

func GetWorkspaces(c *gin.Context) {
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

	results := services.GetWorkspacesWithPagination(filters, opts, query)

	c.JSON(http.StatusOK, models.Response{Data: results})
	return
}

func CreateWorkspace(c *gin.Context) {
	var request models.Workspace

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, models.Response{Data: err.Error()})
		return
	}

	request.Id = uuid.New().String()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	_, err = services.CreateWorkspace(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: request})
	return
}

func GetWorkspaceById(c *gin.Context) {
	id := c.Param("id")

	result := services.GetWorkspace(bson.M{"id": id}, nil)
	if result == nil {
		c.JSON(http.StatusNotFound, models.Result{Data: "Data Not Found"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: result})
}

func UpdateWorkspace(c *gin.Context) {
	id := c.Param("id")

	data := services.GetWorkspace(bson.M{"id": id}, nil)
	if data == nil {
		c.JSON(http.StatusNotFound, models.Result{Data: "Data Not Found"})
		return
	}

	var request models.Workspace

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	_, err = services.UpdateWorkspace(id, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(200, models.Response{Data: request})
}

func DeleteWorkspace(c *gin.Context) {
	id := c.Param("id")

	_, err := services.DeleteWorkspace(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: "Failed Delete Data"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Success"})
}

func GetWorkspaceMembers(c *gin.Context) {
	workspaceId := c.Param("workspaceId")
	workspace := services.GetWorkspace(bson.M{"id": workspaceId}, nil)
	if workspace == nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: "Workspace not found"})
		return
	}

	members := services.GetUsers(bson.M{
		"id": bson.M{
			"$in": workspace.UserIds,
		}}, options.Find().SetProjection(bson.M{"id": 1, "name": 1, "email": 1}))

	c.JSON(http.StatusOK, models.Response{Data: members})
}
