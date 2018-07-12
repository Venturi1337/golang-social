package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login", facebook.handleFacebookLogin)
	http.HandleFunc("/oauth2callback", handleFacebookCallback)
	fmt.Print("Started running on http://localhost:9090\n")
	err := http.ListenAndServeTLS(":9090", "mycert.pem", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
