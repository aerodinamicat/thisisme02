package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aerodinamicat/thisisme02/handlers"
	"github.com/aerodinamicat/thisisme02/middlewares"
	"github.com/aerodinamicat/thisisme02/servers"
	"github.com/gorilla/mux"
)

func setEndPointsHandlers(server *servers.HttpServer, router *mux.Router) {
	router.Use(middlewares.CheckAuthMiddleware(server))

	router.HandleFunc("/signup", handlers.SignUpHandler(server)).Methods(http.MethodPost)
	router.HandleFunc("/login", handlers.LogInHandler(server)).Methods(http.MethodPost)

	router.HandleFunc("/user/changeEmail", handlers.ChangeEmailHandler(server)).Methods(http.MethodPut)
	router.HandleFunc("/user/changePassword", handlers.ChangePasswordHandler(server)).Methods(http.MethodPut)

	// Hay que añadir los endpoints para cambiar el email y cambiar la contraseña.
}

func main() {
	//port := "12345"
	//jwtSecret := "mysecretphrase"
	//dbHost := "localhost"
	//dbSchema := "thisisme"
	//dbUser := "postgres"
	//dbPassword := "mysecretpassword"
	//dbPort := "54321"

	serverConfiguration := &servers.Config{
		Port:   os.Getenv("APP_PORT"),
		Secret: os.Getenv("APP_JWTSECRET"),
	}
	databaseConfiguration := &servers.DBConfig{
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		Host:     os.Getenv("DB_HOST"),
		Schema:   os.Getenv("DB_SCHEMA"),
	}
	httpServer := servers.NewHttpServer(context.Background(), serverConfiguration, databaseConfiguration)
	httpServer.Start(setEndPointsHandlers)
}
