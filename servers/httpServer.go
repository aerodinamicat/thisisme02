package servers

import (
	"context"
	"log"
	"net/http"

	"github.com/aerodinamicat/thisisme02/databases"
	"github.com/gorilla/mux"
)

type Config struct {
	Port   string
	Secret string
}
type DBConfig struct {
	Port     string
	Password string
	User     string
	Host     string
	Schema   string
}

type HttpServer struct {
	Config   *Config
	DBConfig *DBConfig
	Router   *mux.Router
}

func NewHttpServer(ctx context.Context, cfg *Config, dbCfg *DBConfig) *HttpServer {
	server := &HttpServer{
		Config:   cfg,
		DBConfig: dbCfg,
		Router:   mux.NewRouter(),
	}

	return server
}

func (srv *HttpServer) Start(binder func(server *HttpServer, router *mux.Router)) {
	binder(srv, srv.Router)

	dbr, err := databases.NewPostgresImplementation(
		srv.DBConfig.User,
		srv.DBConfig.Password,
		srv.DBConfig.Host,
		srv.DBConfig.Port,
		srv.DBConfig.Schema,
	)
	if err != nil {
		log.Fatalf("Database connection failed: '%v'", err)
	}
	databases.SetDatabaseRepository(dbr)

	log.Printf("Server started and listening for requests on port: %s\n", srv.Config.Port)
	if err := http.ListenAndServe(":"+srv.Config.Port, srv.Router); err != nil {
		log.Fatalf("Error from 'Listen&Serve': '%v'", err)
	}
}
