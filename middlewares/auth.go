package middlewares

import (
	"net/http"
	"strings"

	"github.com/aerodinamicat/thisisme02/handlers"
	"github.com/aerodinamicat/thisisme02/servers"
	"github.com/golang-jwt/jwt"
)

var (
	NO_AUTH_NEEDED_PATHS = []string{
		"/",
		"signup",
		"login",
	}
)

func shouldCheckAuthorization(route string) bool {
	for _, path := range NO_AUTH_NEEDED_PATHS {
		if strings.Contains(route, path) {
			return false
		}
	}
	return true
}

func CheckAuthMiddleware(server *servers.HttpServer) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if shouldCheckAuthorization(request.URL.Path) {
				authTokenString := strings.TrimSpace(request.Header.Get("Authorization"))

				_, err := jwt.ParseWithClaims(authTokenString, &handlers.UserClaim{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(server.Config.Secret), nil
				})
				if err != nil {
					http.Error(writer, err.Error(), http.StatusUnauthorized)
					return
				}
			}
			next.ServeHTTP(writer, request)
		})
	}
}
