package routing

import (
	"github.com/dimitarvalkanov7/chaoscamp-demo/controllers"
	"github.com/gorilla/mux"
)

func Handlers() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	// r.Use(CommonMiddleware)

	// r.HandleFunc("/", controllers.TestAPI).Methods("GET")
	// r.HandleFunc("/api", controllers.TestAPI).Methods("GET")
	r.HandleFunc("/", controllers.Home).Methods("GET")
	r.HandleFunc("/home", controllers.Home).Methods("GET")
	r.HandleFunc("/register", controllers.AddNewUser).Methods("GET")
	r.HandleFunc("/register", controllers.Register).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("GET")
	r.HandleFunc("/auth", controllers.Authenticate).Methods("POST")

	// Auth route
	// s := r.PathPrefix("/auth").Subrouter()
	// s.Use(auth.JwtVerify)
	// s.HandleFunc("/user", controllers.FetchUsers).Methods("GET")
	// s.HandleFunc("/user/{id}", controllers.GetUser).Methods("GET")
	// s.HandleFunc("/user/{id}", controllers.UpdateUser).Methods("PUT")
	// s.HandleFunc("/user/{id}", controllers.DeleteUser).Methods("DELETE")
	return r
}
