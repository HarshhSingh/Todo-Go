package middlewares

import (
	"fmt"
	"main/utils"
	"net/http"
)

func JWTAuthorisation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("JWT Authorisation")
		headerToken := r.Header.Get("Authorization")
		fmt.Printf("headerToken: %v\n", headerToken)
		if headerToken == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		headerToken = headerToken[len("Bearer "):]

		err := utils.VerifyToken(headerToken)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
