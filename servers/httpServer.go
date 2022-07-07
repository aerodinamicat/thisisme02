package servers

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	Config *Config
	Router *mux.Router
}

type Config struct {
	Port            string
	JWTSecretPhrase string
}

func NewHTTPServer(ctx context.Context, port string, jwtSecretPhrase string) *HTTPServer {
	config := &Config{
		Port:            port,
		JWTSecretPhrase: jwtSecretPhrase,
	}

	return &HTTPServer{
		Config: config,
		Router: mux.NewRouter(),
	}
}

func (srv *HTTPServer) stablishHandlers() {
	//srv.Router.Use(middleware.CheckAuthNeeded(srv))

	//srv.Router.HandlerFunc("/signup", handlers.SignUpHandler(srv).Methods(http.MethodPost))
}

func (srv *HTTPServer) Start() {
	var err error

	srv.stablishHandlers()

	log.Printf("Starting server on port: %s\n", srv.Config.Port)
	if err = http.ListenAndServe(srv.Config.Port, srv.Router); err != nil {
		log.Fatalf("Error from 'Listen&Serve: '%v'", err)
	}
}
