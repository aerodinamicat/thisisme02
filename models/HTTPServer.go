package models

import database "github.com/aerodinamicat/thisisme02/database"

type HTTPServer struct {
	Config   *Config
	DBServer *database.DatabaseRepository
}

type Config struct {
	Port string
}
