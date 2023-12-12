package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"log"
	"net/http"
	"os"
)

var mySigningKey = []byte(os.Getenv("SECRET_KEY"))

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "info")
}

func fetchData(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("http://localhost:8080")

	if err != nil {
		http.Error(w, "Error fetching token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	tokenBytes, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Error reading token response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString := string(tokenBytes)

	isAuthorized(homePage, tokenString)(w, r)
}

func isAuthorized(endpoint http.HandlerFunc, tokenString string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("not valid")
			}
			return mySigningKey, nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if token.Valid {
			endpoint(w, r)
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}
	}
}

func handleRequests() {
	http.HandleFunc("/", fetchData)
	log.Fatal(http.ListenAndServe(":9001", nil))
}

func main() {
	fmt.Printf("Starting the server at port 9001\n")
	handleRequests()
}
