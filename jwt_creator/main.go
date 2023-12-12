package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte(os.Getenv("SECRET_KEY"))

func GetJWT() (string, error){
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer: "jwtgo.io",
		ExpiresAt: time.Now().Add(time.Hour*1 ).Unix(),
		Audience: "billing.jwtcreator.com",
	})

	token, err := claims.SignedString(mySigningKey)

	if err!=nil{
		fmt.Errorf("Something went wrong: %s", err.Error())
		return "", err
	}
	return token, nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	validToken, err := GetJWT()
	fmt.Println(validToken)
	if err != nil {
		fmt.Println("Failed to create token")
	}
	fmt.Fprintf(w, string(validToken))
}

func handleRequests() {
	http.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
