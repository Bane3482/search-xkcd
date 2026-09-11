package middleware

import (
	"net/http"
	"strings"
)

type TokenVerifier interface {
	Verify(token string) error
}

func Auth(next http.HandlerFunc, verifier TokenVerifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		tokenInfo := strings.Fields(authHeader)

		if len(tokenInfo) != 2 || tokenInfo[0] != "Token" {
			http.Error(w, "incorrect token format", http.StatusUnauthorized)
			return
		}

		token := tokenInfo[1]

		if err := verifier.Verify(token); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
