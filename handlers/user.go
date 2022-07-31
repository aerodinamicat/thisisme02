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
	HASH_COST            = 8
	EXPIRE_TIME          = 1 * time.Hour * 24
	HEADER_AUTHORIZATION = "Authorization"
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
	Result bool `json:"result"`
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

type PropertyChangesListRequest struct {
	PageInfo models.PageInfo `json:"pageInfo"`
	Name     string          `json:"propertyName"`
}
type PropertyChangesListResponse struct {
	PageInfo        *models.PageInfo         `json:"pageInfo"`
	PropertyChanges []*models.PropertyChange `json:"propertyName"`
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
		//* Preparamos la petición y la recibimos.
		var decodedRequest = new(SignUpRequest)
		if err := json.NewDecoder(request.Body).Decode(&decodedRequest); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		//* Generamos un nuevo 'id' aleatorio.
		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Ciframos la 'password'
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(decodedRequest.Password), HASH_COST)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Guardamos el momento en el que se produjo el registro.
		currentTime := time.Now()

		//* Instanciamos un 'user' con todas sus propiedades.
		var user = models.User{
			Id:        id.String(),
			Email:     decodedRequest.Email,
			Password:  string(hashedPassword),
			CreatedAt: currentTime,
			CreatedBy: id.String(),
			UpdatedAt: currentTime,
			UpdatedBy: id.String(),
		}

		//* Guardamos el 'user' en la DB.
		if err := databases.InsertUser(request.Context(), &user); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Instanciamos los cambios de propiedades pertinentes:
		//* Para 'email'.
		propertyChange := &models.PropertyChange{
			UserId:    user.Id,
			Name:      "email",
			From:      "",
			To:        user.Email,
			CreatedAt: currentTime,
			CreatedBy: user.Id,
		}
		//* Lo guardamos en ls DB.
		if err := databases.InsertPropertyChangeLog(request.Context(), propertyChange); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Para 'password'.
		propertyChange = &models.PropertyChange{
			UserId:    user.Id,
			Name:      "password",
			From:      "",
			To:        user.Password,
			CreatedAt: currentTime,
			CreatedBy: user.Id,
		}
		//* Lo guardamos.
		if err := databases.InsertPropertyChangeLog(request.Context(), propertyChange); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Preparamos la respuesta y la enviamos.
		notEncodedResponse := SignUpResponse{
			Result: true,
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(notEncodedResponse)
	}
}
func LogInHandler(server *servers.HttpServer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//* Preparamos la petición y la recibimos.
		var decodedRequest = new(LogInRequest)
		if err := json.NewDecoder(request.Body).Decode(&decodedRequest); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		//* Solicitamos a DB un 'user' con los mismos datos que los facilitados.
		user, err := databases.GetUserByEmail(request.Context(), decodedRequest.Email)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Si 'user' devuelto es nulo, es decir no existe en DB, respondemos 'No autorizado'.
		if user == nil {
			http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		//* Si 'user' existe, comparamos su 'password' con la facilitada en la petición.
		//* Si 'error' devuelto no es nulo, es decir las 'password' no coinciden, respondemos 'No autorizado'.
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(decodedRequest.Password)); err != nil {
			http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		//* Si las 'password' coinciden, generamos un token de autorización.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, newUserClaim(user.Id))
		authorizationToken, err := token.SignedString([]byte(server.Config.Secret))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Preparamos la respuesta y la enviamos.
		notEncodedResponse := &LogInResponse{
			Authorization: authorizationToken,
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set(HEADER_AUTHORIZATION, authorizationToken)
		json.NewEncoder(writer).Encode(notEncodedResponse)
	}
}
func ChangeEmailHandler(server *servers.HttpServer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//* Esta función requiere autenticación de usuario.
		//* Obtenemos el 'token' de las 'headers' de la petición. Se ubican en 'Authorization'.
		authorizationToken := strings.TrimSpace(request.Header.Get(HEADER_AUTHORIZATION))
		token, err := jwt.ParseWithClaims(authorizationToken, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(server.Config.Secret), nil
		})
		//* Si no podemos decodificarla, devolvemos 'No autorizado'.
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		//* En caso afirmativo, obtenemos un 'id' de 'user' en ella.
		claims, ok := token.Claims.(*UserClaim)
		if !ok || !token.Valid {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Solicitamos un 'user' a DB.
		user, err := databases.GetUserById(request.Context(), claims.UserId)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Preparamos la petición y la recibimos
		var decodedRequest = new(ChangeEmailRequest)
		if err := json.NewDecoder(request.Body).Decode(&decodedRequest); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		//* Guardamos el momento en el que se produjo el registro.
		currentTime := time.Now()

		//* Instanciamos un 'propertyChange'.
		propertyChange := &models.PropertyChange{
			UserId:    user.Id,
			Name:      "email",
			From:      user.Email,
			To:        decodedRequest.NewEmail,
			CreatedAt: currentTime,
			CreatedBy: user.Id,
		}

		//* Realizamos el cambio en 'user'.
		user.Email = decodedRequest.NewEmail
		user.UpdatedAt = currentTime

		//* Guardamos 'user' en DB.
		if err := databases.UpdateUser(request.Context(), user); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Guardamos 'propertyChange' en DB.
		if err := databases.InsertPropertyChangeLog(request.Context(), propertyChange); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Preparamos la respuesta y la enviamos.
		notEncodedResponse := ChangeEmailResponse{
			Result: true,
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(notEncodedResponse)
	}
}
func ChangePasswordHandler(server *servers.HttpServer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//* Esta función requiere autenticación de usuario.
		//* Obtenemos el 'token' de las 'headers' de la petición. Se ubican en 'Authorization'.
		authorizationToken := strings.TrimSpace(request.Header.Get(HEADER_AUTHORIZATION))
		token, err := jwt.ParseWithClaims(authorizationToken, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(server.Config.Secret), nil
		})
		//* Si no podemos decodificarla, devolvemos 'No autorizado'.
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		//* En caso afirmativo, obtenemos un 'id' de 'user' en ella.
		claims, ok := token.Claims.(*UserClaim)
		if !ok || !token.Valid {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Solicitamos un 'user' a DB.
		user, err := databases.GetUserById(request.Context(), claims.UserId)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Preparamos la petición y la recibimos
		var decodedRequest = new(ChangePasswordRequest)
		if err := json.NewDecoder(request.Body).Decode(&decodedRequest); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		//* Comparamos ambas 'password'. Si no coinciden, devolvemos 'No autorizado'.
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(decodedRequest.CurrentPassword)); err != nil {
			http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		//* Ciframos la nueva 'password'.
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(decodedRequest.NewPassword), HASH_COST)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Guardamos el momento en el que se produjo el registro.
		currentTime := time.Now()

		//* Instanciamos un 'propertyChange'.
		propertyChange := &models.PropertyChange{
			UserId:    user.Id,
			Name:      "password",
			From:      user.Password,
			To:        string(hashedPassword),
			CreatedAt: currentTime,
			CreatedBy: user.Id,
		}

		//* Realizamos el cambio en 'user'.
		user.Password = string(hashedPassword)
		user.UpdatedAt = currentTime

		//* Guardamos 'user' en DB.
		if err := databases.UpdateUser(request.Context(), user); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Guardamos 'propertyChange' en DB.
		if err := databases.InsertPropertyChangeLog(request.Context(), propertyChange); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Preparamos la respuesta y la enviamos.
		notEncodedResponse := ChangePasswordResponse{
			Result: true,
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(notEncodedResponse)
	}
}
func GetPropertyChangesHandler(server *servers.HttpServer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//* Esta función requiere autenticación de usuario.
		//* Obtenemos el 'token' de las 'headers' de la petición. Se ubican en 'Authorization'.
		authorizationToken := strings.TrimSpace(request.Header.Get(HEADER_AUTHORIZATION))
		token, err := jwt.ParseWithClaims(authorizationToken, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(server.Config.Secret), nil
		})
		//* Si no podemos decodificarla, devolvemos 'No autorizado'.
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		//* En caso afirmativo, obtenemos un 'id' de 'user' en ella.
		claims, ok := token.Claims.(*UserClaim)
		if !ok || !token.Valid {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Preparamos la petición y la recibimos
		var decodedRequest = new(PropertyChangesListRequest)
		if err := json.NewDecoder(request.Body).Decode(&decodedRequest); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		//* Solicitamos un listado 'propertyChange'  a DB.
		propertyChanges, pageInfo, err := databases.ListPropertyChangesByUserIdAndName(request.Context(), claims.UserId, decodedRequest.Name, &decodedRequest.PageInfo)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//* Preparamos la respuesta y la enviamos.
		notEncodedResponse := &PropertyChangesListResponse{
			PageInfo:        pageInfo,
			PropertyChanges: propertyChanges,
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(notEncodedResponse)
	}
}
