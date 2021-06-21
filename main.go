package main

import (
	"fmt"
	"jwt-demo/handler"
	"log"
	"net/http"
)

/*
POST http://localhost:8000/auth
{"username":"user1","password":"password1"}
*/
func main() {
	http.HandleFunc("/signin", handler.Signin)
	http.HandleFunc("/welcome", handler.Welcome)
	http.HandleFunc("/refresh", handler.Refresh)

	//auth token
	http.HandleFunc("/auth", handler.Auth)
	//validate token
	http.HandleFunc("/validate", handler.Validate)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", handler.Port), nil); err != nil {
		log.Fatalln(err)
	}
}
