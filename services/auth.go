package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"kickof/database"
	"kickof/models"
	"kickof/utils"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func GenerateToken(email string) (*string, error) {
	expire := time.Now().Add(24 * time.Hour)

	claims := models.Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	if err == nil {
		return &tokenString, nil
	}

	return nil, errors.New(err.Error())
}

func CheckHeader(r *http.Request) (string, error) {
	header := r.Header["Authorization"]

	if header == nil {
		return "", errors.New("unauthorized")
	}

	split := strings.Split(header[0], " ")
	if len(split) != 2 || strings.ToLower(split[0]) != "bearer" {
		return "", errors.New("unauthorized")
	}

	return split[1], nil
}

func VerifyToken(tokenString string) (string, error) {
	claims := &models.Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		log.Println(err)
		return "", errors.New("Token invalid")
	}

	return claims.Email, nil
}

func SignIn(params models.Login) (*string, error) {
	user := GetUser(bson.M{"email": params.Email}, nil)

	if user == nil {
		return nil, errors.New("User not found")
	}

	pass := utils.ComparePassword(user.Password, []byte(params.Password))

	if !pass {
		return nil, errors.New("Password does not match")
	}

	token, err := GenerateToken(params.Email)

	if err == nil {
		data := user
		data.LastActive = time.Now()

		_, err = database.UpdateOne(UserCollection, bson.M{"email": params.Email}, data)
		if err != nil {
			return nil, err
		}

		return token, nil
	}

	return nil, errors.New(err.Error())
}

func Register(params models.Register, url string) (*string, error) {
	email := GetUser(bson.M{"email": params.Email}, nil)

	if email != nil {
		return nil, errors.New("Email already exists")
	}

	request := models.User{}
	request.Id = uuid.New().String()
	request.Email = params.Email
	request.Password = utils.HashAndSalt(params.Password)
	request.Name = params.Name
	request.CreatedAt = time.Now()
	request.LastActive = time.Now()

	_, e := database.InsertOne(UserCollection, request)

	if e != nil {
		return nil, e
	}

	token, err := GenerateToken(params.Email)

	if err != nil {
		return nil, err
	}

	tokenValue := *token

	data := models.VerificationMail{
		Name: request.Name,
		Link: url + tokenValue,
	}

	_, err = utils.SendEmailVerification("verification", request.Email, data)

	if err != nil {
		log.Println("Failed to send verification email")
	}

	return token, nil
}

func Activate(token string) (bool, error) {
	email, err := VerifyToken(token)

	if err != nil {
		return false, err
	}

	filter := bson.M{"email": email}

	var user models.User
	cursor := database.FindOne(UserCollection, filter, nil)

	err = cursor.Decode(&user)
	if err != nil {
		return false, err
	}

	user.Active = true
	user.Status = true

	_, err = database.UpdateOne(UserCollection, filter, user)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ForgotPassword(email string, url string) (bool, error) {
	user := GetUser(bson.M{"email": email}, nil)
	if user == nil {
		return false, errors.New("user not found")
	}

	token, err := GenerateToken(email)

	if err != nil {
		return false, err
	}

	tokenValue := *token

	data := models.VerificationMail{
		Name: user.Name,
		Link: url + tokenValue,
	}

	res, err := utils.SendEmailVerification("forgot-password", user.Email, data)

	if err != nil {
		return false, err
	}

	return res, nil
}

func UpdatePassword(token string, password string) (bool, error) {
	email, err := VerifyToken(token)

	if err != nil {
		return false, err
	}

	filter := bson.M{"email": email}

	var user models.User
	cursor := database.FindOne(UserCollection, filter, nil)

	err = cursor.Decode(&user)
	if err != nil {
		return false, err
	}

	user.Password = utils.HashAndSalt(password)

	_, err = database.UpdateOne(UserCollection, filter, user)
	if err != nil {
		return false, err
	}

	return true, nil
}
