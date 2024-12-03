package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kickof/database"
	"kickof/models"
	"net/http"
)

const UserCollection = "users"

func GetUsers(filters bson.M, opt *options.FindOptions) []models.User {
	results := make([]models.User, 0)

	cursor := database.Find(UserCollection, filters, opt)

	if cursor == nil {
		return results
	}
	for cursor.Next(context.Background()) {
		var user models.User

		if cursor.Decode(&user) == nil {
			results = append(results, user)
		}
	}
	return results
}

func GetUsersWithPagination(filters bson.M, opt *options.FindOptions, query models.Query) models.Result {
	results := GetUsers(filters, opt)

	count := database.Count(UserCollection, filters)

	pagination := query.GetPagination(count)

	result := models.Result{
		Data:       results,
		Pagination: pagination,
		Query:      query,
	}

	return result
}

func CreateUser(user models.User) (bool, error) {
	_, err := database.InsertOne(UserCollection, user)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetUser(filter bson.M, opts *options.FindOneOptions) *models.User {

	cursor := database.FindOne(UserCollection, filter, opts)

	if cursor == nil {
		return nil
	}

	var user models.User

	err := cursor.Decode(&user)
	if err != nil {
		return nil
	}

	return &user
}

func GetUserByEmail(email string) *models.User {
	cursor := database.FindOne(UserCollection, bson.M{"email": email}, nil)
	if cursor == nil {
		return nil
	}

	var user models.User
	err := cursor.Decode(&user)
	if err != nil {
		return nil
	}

	user.Password = ""

	return &user
}

func UpdateUser(id string, user models.User) (*mongo.UpdateResult, error) {
	filters := bson.M{"id": id}

	res, err := database.UpdateOne(UserCollection, filters, user)

	if res == nil {
		return nil, err
	}

	return res, nil
}

func DeleteUser(id string) (*mongo.DeleteResult, error) {
	filter := bson.M{"id": id}

	res, err := database.DeleteOne(UserCollection, filter)

	if res == nil {
		return nil, err
	}

	return res, nil
}

func GetCurrentUser(r *http.Request) *models.User {
	token, e := CheckHeader(r)

	if e != nil {
		return nil
	}

	email, err := VerifyToken(token)

	if err != nil {
		return nil
	}

	user := GetUser(bson.M{"email": email}, options.FindOne().SetProjection(bson.D{{"password", 0}}))

	return user
}
