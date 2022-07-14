package main

import (
	"context"
	"net/http"

	"github.com/aerodinamicat/thisisme02/handlers"
	"github.com/aerodinamicat/thisisme02/servers"
	"github.com/gorilla/mux"
)

func StablishEndPointsHandlers(server servers.HttpServer, router *mux.Router) {
	//router.Use(middlewares.CheckAuthMiddleware(server))

	router.HandleFunc("/signup", handlers.SignUpHandler(server)).Methods(http.MethodPost)
	router.HandleFunc("/login", handlers.LogInHandler(server)).Methods(http.MethodPost)

}

func main() {
	port := "12345"
	jwtSecret := "mysecretphrase"
	dbHost := "localhost"
	dbSchema := "thisisme"
	dbUser := "postgres"
	dbPassword := "mysecretpassword"
	dbPort := "54321"

	cfg := &servers.Config{
		Port:   port,
		Secret: jwtSecret,
	}
	dbCfg := &servers.DBConfig{
		Port:     dbPort,
		Password: dbPassword,
		User:     dbUser,
		Host:     dbHost,
		Schema:   dbSchema,
	}
	httpServer := servers.NewHttpServer(context.Background(), cfg, dbCfg)

	httpServer.Start(StablishEndPointsHandlers)
}
