package routes

import "github.com/gorilla/mux"

var RegisteredBookStoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/book/", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/book/", controllers.GetBook).Methods("GET")
	router.HandleFunc("/book/{bookId}", controllers.GetBookById).Methods("GET")
	router.HandleFunc("/book/", controllers.UpdateBook).Methods("PUT")
	router.HandleFunc("/book/", controllers.DeleteBook).Methods("DELETE")
}
