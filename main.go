package main

import (
	"fmt"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

var privateHtml = `
Hello %s <br />
<a href="/login?logout=true">logout</a>
`

var publicHtml = `
<a href="/login">login</a>
`

func main() {
	secret := os.Getenv("LOGINSRV_JWT_SECRET")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html; charset=utf8")
		if c, err := r.Cookie("jwt_token"); err == nil {
			token, err := jwt.Parse(c.Value, func(*jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				fmt.Fprintf(w, privateHtml, claims["sub"])
				return
			}
		}
		fmt.Fprintln(w, publicHtml)
	})
	http.ListenAndServe(":8888", nil)
}
