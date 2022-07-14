package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aerodinamicat/thisisme02/databases"
	"github.com/aerodinamicat/thisisme02/models"
	"github.com/aerodinamicat/thisisme02/servers"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	HASH_COST   = 8
	EXPIRE_TIME = 1 * time.Hour * 24
)

type UserClaim struct {
	UserId string
	jwt.StandardClaims
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SignUpResponse struct {
	Id string `json:"id"`
}
type LogInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LogInResponse struct {
	AuthorizationToken string `json:"authorization"`
}

func newUserClaim(userId string) UserClaim {
	return UserClaim{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(EXPIRE_TIME).Unix(),
		},
	}
}

func SignUpHandler(server servers.HttpServer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var decodedRequest SignUpRequest
		var notEncodedResponse SignUpResponse

		decodedRequest = SignUpRequest{}
		if err := json.NewDecoder(request.Body).Decode(&decodedRequest); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(decodedRequest.Password), HASH_COST)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		var user = models.User{
			Id:       id.String(),
			Email:    decodedRequest.Email,
			Password: string(hashedPassword),
		}

		if err := databases.InsertUser(request.Context(), &user); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		notEncodedResponse = SignUpResponse{
			Id: user.Id,
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(notEncodedResponse)
	}
}
func LogInHandler(server servers.HttpServer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var decodedRequest LogInRequest
		var notEncodedResponse LogInResponse
		var user *models.User

		if err := json.NewDecoder(request.Body).Decode(&decodedRequest); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := databases.GetUserByEmail(request.Context(), decodedRequest.Email)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(decodedRequest.Password)); err != nil {
			http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, newUserClaim(user.Id))
		authorizationToken, err := token.SignedString([]byte(server.Config().Secret))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		notEncodedResponse = LogInResponse{
			AuthorizationToken: authorizationToken,
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(notEncodedResponse)
	}
}
