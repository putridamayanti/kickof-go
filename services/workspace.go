package services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kickof/database"
	"kickof/models"
)

const WorkspaceCollection = "workspaces"

func GetWorkspaceMembers(userIds []string) []models.User {
	result := GetUsers(bson.M{"id": bson.M{"$in": userIds}}, options.Find().SetProjection(bson.D{{"password", 0}}))
	return result
}

func GetWorkspaces(filters bson.M, opt *options.FindOptions) []models.Workspace {
	results := make([]models.Workspace, 0)

	cursor := database.Find(WorkspaceCollection, filters, opt)
	if cursor == nil {
		return results
	}
	for cursor.Next(context.Background()) {
		var data models.Workspace
		err := cursor.Decode(&data)
		if err == nil {
			data.Members = GetWorkspaceMembers(data.UserIds)
			results = append(results, data)
		}

	}

	return results
}

func GetWorkspacesWithPagination(filters bson.M, opt *options.FindOptions, query models.Query) models.Result {
	results := GetWorkspaces(filters, opt)

	count := database.Count(WorkspaceCollection, filters)

	pagination := query.GetPagination(count)

	result := models.Result{
		Data:       results,
		Pagination: pagination,
		Query:      query,
	}

	return result
}

func CreateWorkspace(Workspace models.Workspace) (bool, error) {
	_, err := database.InsertOne(WorkspaceCollection, Workspace)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetWorkspace(filter bson.M, opts *options.FindOneOptions) *models.Workspace {
	var data models.Workspace
	err := database.FindOne(WorkspaceCollection, filter, opts).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return nil
	}
	return &data
}

func UpdateWorkspace(id string, Workspace interface{}) (*mongo.UpdateResult, error) {
	filters := bson.M{"id": id}

	res, err := database.UpdateOne(WorkspaceCollection, filters, Workspace)

	if res == nil {
		return nil, err
	}

	return res, nil
}

func DeleteWorkspace(id string) (*mongo.DeleteResult, error) {
	filter := bson.M{"id": id}

	res, err := database.DeleteOne(WorkspaceCollection, filter)

	if res == nil {
		return nil, err
	}

	return res, nil
}
