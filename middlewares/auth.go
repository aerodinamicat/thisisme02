package middlewares

import (
	"net/http"
	"strings"

	"github.com/aerodinamicat/thisisme02/handlers"
	"github.com/aerodinamicat/thisisme02/servers"
	"github.com/golang-jwt/jwt"
)

func CheckAuthMiddleware(server servers.HttpServer) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			authTokenString := strings.TrimSpace(request.Header.Get("Authorization"))

			_, err := jwt.ParseWithClaims(authTokenString, &handlers.UserClaim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(server.Config().Secret), nil
			})
			if err != nil {
				http.Error(writer, err.Error(), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}
