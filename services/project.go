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

const ProjectCollection = "projects"

func GetProjectMembers(userIds []string) []models.User {
	result := GetUsers(bson.M{"id": bson.M{"$in": userIds}}, options.Find().SetProjection(bson.D{{"password", 0}}))
	return result
}

func GetProjects(filters bson.M, opt *options.FindOptions) []models.Project {
	results := make([]models.Project, 0)

	cursor := database.Find(ProjectCollection, filters, opt)
	if cursor == nil {
		return results
	}
	for cursor.Next(context.Background()) {
		var data models.Project
		err := cursor.Decode(&data)
		if err == nil {
			data.Members = GetProjectMembers(data.UserIds)
			//if data.WorkspaceId != "" {
			//	workspace := GetWorkspace(bson.M{"id": data.WorkspaceId}, nil)
			//	if workspace != nil {
			//		data.Workspace = *workspace
			//	}
			//}

			results = append(results, data)
		}

	}

	return results
}

func GetProjectsWithPagination(filters bson.M, opt *options.FindOptions, query models.Query) models.Result {
	results := GetProjects(filters, opt)

	count := database.Count(ProjectCollection, filters)

	pagination := query.GetPagination(count)

	result := models.Result{
		Data:       results,
		Pagination: pagination,
		Query:      query,
	}

	return result
}

func CreateProject(Project models.Project) (bool, error) {
	_, err := database.InsertOne(ProjectCollection, Project)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetProject(filter bson.M, opts *options.FindOneOptions) *models.Project {
	var data models.Project
	err := database.FindOne(ProjectCollection, filter, opts).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return nil
	}
	return &data
}

func UpdateProject(id string, Project interface{}) (*mongo.UpdateResult, error) {
	filters := bson.M{"id": id}

	res, err := database.UpdateOne(ProjectCollection, filters, Project)

	if res == nil {
		return nil, err
	}

	return res, nil
}

func DeleteProject(id string) (*mongo.DeleteResult, error) {
	filter := bson.M{"id": id}

	res, err := database.DeleteOne(ProjectCollection, filter)

	if res == nil {
		return nil, err
	}

	return res, nil
}
