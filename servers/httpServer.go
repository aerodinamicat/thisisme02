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

type HttpServer interface {
	Config() *Config
}

type Broker struct {
	config   *Config
	dbConfig *DBConfig
	router   *mux.Router
}

func (b *Broker) Config() *Config {
	return b.config
}

func NewHttpServer(ctx context.Context, cfg *Config, dbCfg *DBConfig) *Broker {
	broker := &Broker{
		config:   cfg,
		dbConfig: dbCfg,
		router:   mux.NewRouter(),
	}

	return broker
}

func (b *Broker) Start(binder func(server HttpServer, router *mux.Router)) {
	binder(b, b.router)

	dbr, err := databases.NewPostgresImplementation(b.dbConfig.User, b.dbConfig.Password, b.dbConfig.Host, b.dbConfig.Port, b.dbConfig.Schema)
	if err != nil {
		log.Fatalf("Database connection failed: '%v'", err)
	}
	databases.SetDatabaseRepository(dbr)

	//log.Printf("Starting server on port: %s\n", b.config.Port)
	if err := http.ListenAndServe(":"+b.config.Port, b.router); err != nil {
		log.Fatalf("Error from 'Listen&Serve: '%v'", err)
	}
}
