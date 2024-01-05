package main

import (
	"fmt"
	"net/http"
)

func main() {
	r := router.Router()
	fmt.Println("Starting the server at port 8000..")

	http.ListenAndServe(":8000", r)

}
