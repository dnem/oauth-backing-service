package server

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

func helloHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 {
			w.Header().Set("Content-Type", "text/html;charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			buf := bytes.NewBufferString(`INVALID REQUEST`)
			w.Write(buf.Bytes())
		}

		token := authHeaderParts[1]

		accessToken, err := parseToken(token)
		if err != nil {
			fmt.Printf("Error Parsing Token: %s\n", err)
		}

		user := accessToken.Claims["user_name"].(string)
		if len(user) == 0 {
			user = "Unknown User"
		}

		if hasScope(accessToken, "test.access", "test.admin") {
			w.Header().Set("Content-Type", "text/html;charset=utf-8")
			buf := bytes.NewBufferString(user + " ACCESSED BACKING SERVICE")
			w.Write(buf.Bytes())
		} else {
			w.Header().Set("Content-Type", "text/html;charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			buf := bytes.NewBufferString(user + " USER DOES NOT HAVE THE APPROPRIATE PRIVILEGE")
			w.Write(buf.Bytes())
		}
	}
}
