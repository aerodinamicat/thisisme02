package servers

import (
	"github.com/aerodinamicat/thisisme02/databases"
	"github.com/gorilla/mux"
)

type HTTPServer struct {
	Config   *Config
	Router   *mux.Router
	DBServer *databases.DatabaseRepository
}

type Config struct {
	Port string
}
