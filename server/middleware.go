package server

import (
	"net/http"
	"strings"
)

func (mem *Memory) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Split(r.Header["Authorization"][0], " ")[1] != mem.ApiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
