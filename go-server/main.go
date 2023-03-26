package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileServer)

	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)
	fmt.Println("Server is running at port 4000")

	log.Fatal(http.ListenAndServe(":4000", nil))
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "This is my form handler\n")
	firstname := r.FormValue("fname")
	lastname := r.FormValue("lname")

	fmt.Fprintf(w, "First Name : %s \nLast Name : %s", firstname, lastname)

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
	}

	if r.Method != "GET" {
		http.Error(w, "Method not found", http.StatusNotFound)
	}
	fmt.Fprintf(w, "Hello, It is my first web server")
}
