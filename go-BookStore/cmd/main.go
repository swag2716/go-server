package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swapnika/go-bookstore/pkg/routes"
)

func main() {
	r := mux.NewRouter()
	routes.RegisteredBookStoreRoutes(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", r))
}
