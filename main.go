package main

import (
	"log"
	"net/http"

	"github.com/Harits514/AdroadyTes/app"
)

func main() {
	log.Println("server stated at localhost/8080")
	http.Handle("/", http.FileServer(http.Dir("/")))
	http.ListenAndServe(":8080", app.Router.Route)
}
