package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"kickof/models"
	"kickof/services"
	"log"
	"net/http"
	"slices"
	"time"
)

func GetTasks(c *gin.Context) {
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

	if query.Assigned == "true" {
		filters["assigneeIds"] = bson.M{
			"$exists": true,
			"$not": bson.M{
				"$size": 0,
			},
		}
	}

	if query.Completed == "true" {
		state := services.GetState(bson.M{"name": "Done"}, nil)
		if state != nil {
			filters["stateId"] = state.Id
		}
	}

	opts := query.GetOptions()

	results := services.GetTasksWithPagination(filters, opts, query)

	c.JSON(http.StatusOK, models.Response{Data: results})
	return
}

func CreateTask(c *gin.Context) {
	var request models.Task

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, models.Response{Data: err.Error()})
		return
	}

	request.Id = uuid.New().String()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	for _, val := range request.Labels {
		if val.Id == "" {
			val.Id = uuid.New().String()
			val.WorkspaceId = request.WorkspaceId
			val.ProjectId = request.ProjectId
			val.CreatedAt = time.Now()
			val.UpdatedAt = time.Now()
			_, err = services.CreateTaskLabel(val)
			if err != nil {
				log.Println("Error create task label", err.Error())
			}
		}
		request.LabelIds = append(request.LabelIds, val.Id)
	}

	for _, val := range request.Assignees {
		request.AssigneeIds = append(request.AssigneeIds, val.Id)
	}

	_, err = services.CreateTask(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: request})
	return
}

func GetTaskById(c *gin.Context) {
	id := c.Param("id")

	result := services.GetTask(bson.M{"id": id}, nil)
	if result == nil {
		c.JSON(http.StatusNotFound, models.Result{Data: "Data Not Found"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: result})
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")

	data := services.GetTask(bson.M{"id": id}, nil)
	if data == nil {
		c.JSON(http.StatusNotFound, models.Result{Data: "Data Not Found"})
		return
	}

	var request models.Task
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	for _, val := range request.Labels {
		if val.Id == "" {
			val.Id = uuid.New().String()
			val.WorkspaceId = request.WorkspaceId
			val.ProjectId = request.ProjectId
			val.CreatedAt = time.Now()
			val.UpdatedAt = time.Now()
			_, err = services.CreateTaskLabel(val)
			if err != nil {
				log.Println("Error create task label", err.Error())
			}
		}

		request.LabelIds = append(request.LabelIds, val.Id)
	}

	for _, val := range request.Assignees {
		request.AssigneeIds = append(request.AssigneeIds, val.Id)
	}

	request.UpdatedAt = time.Now()
	slices.Sort(request.LabelIds)
	slices.Compact(request.LabelIds)
	slices.Sort(request.AssigneeIds)
	slices.Compact(request.AssigneeIds)
	_, err = services.UpdateTask(id, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(200, models.Response{Data: request})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	_, err := services.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: "Failed Delete Data"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Success"})
}
