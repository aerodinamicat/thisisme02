package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
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
	Authorization string `json:"authorization"`
}

type ChangeEmailRequest struct {
	NewEmail string `json:"newEmail"`
}
type ChangeEmailResponse struct {
	Result bool `json:"result"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
type ChangePasswordResponse struct {
	Result bool `json:"result"`
}

func newUserClaim(userId string) UserClaim {
	return UserClaim{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(EXPIRE_TIME).Unix(),
		},
	}
}

func SignUpHandler(server *servers.HttpServer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var decodedRequest = new(SignUpRequest)
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

		notEncodedResponse := SignUpResponse{
			Id: user.Id,
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(notEncodedResponse)
	}
}
func LogInHandler(server *servers.HttpServer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var decodedRequest = new(LogInRequest)
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
		authorizationToken, err := token.SignedString([]byte(server.Config.Secret))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		notEncodedResponse := &LogInResponse{
			Authorization: authorizationToken,
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(notEncodedResponse)
	}
}
func ChangeEmailHandler(server *servers.HttpServer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		authorizationToken := strings.TrimSpace(request.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(authorizationToken, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(server.Config.Secret), nil
		})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*UserClaim)
		if !ok || !token.Valid {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := databases.GetUserById(request.Context(), claims.UserId)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		var decodedRequest = new(ChangeEmailRequest)
		if err := json.NewDecoder(request.Body).Decode(&decodedRequest); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		user.Email = decodedRequest.NewEmail

		if err := databases.UpdateUser(request.Context(), user); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		notEncodedResponse := ChangeEmailResponse{
			Result: true,
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(notEncodedResponse)
	}
}
func ChangePasswordHandler(server *servers.HttpServer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		authorizationToken := strings.TrimSpace(request.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(authorizationToken, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(server.Config.Secret), nil
		})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*UserClaim)
		if !ok || !token.Valid {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := databases.GetUserById(request.Context(), claims.UserId)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		var decodedRequest = new(ChangePasswordRequest)
		if err := json.NewDecoder(request.Body).Decode(&decodedRequest); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(decodedRequest.CurrentPassword)); err != nil {
			http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(decodedRequest.NewPassword), HASH_COST)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		user.Password = string(hashedPassword)

		if err := databases.UpdateUser(request.Context(), user); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		notEncodedResponse := ChangePasswordResponse{
			Result: true,
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(notEncodedResponse)
	}
}
