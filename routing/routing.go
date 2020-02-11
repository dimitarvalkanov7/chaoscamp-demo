package routing

import (
	"github.com/dimitarvalkanov7/chaoscamp-demo/controllers"
	"github.com/dimitarvalkanov7/chaoscamp-demo/encryption"
	"github.com/gorilla/mux"
	"net/http"
)

func Handlers() *mux.Router {
	const STATIC_DIR = "/static/"
	r := mux.NewRouter().StrictSlash(true)
	r.
		PathPrefix(STATIC_DIR).
		Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))
	r.HandleFunc("/", controllers.Login).Methods(http.MethodGet)
	r.HandleFunc("/login", controllers.Login).Methods(http.MethodGet)
	r.HandleFunc("/login", controllers.Authenticate).Methods(http.MethodPost)

	s := r.PathPrefix("/auth").Subrouter()
	s.Use(AuthenticationMiddleware)
	s.HandleFunc("/home", controllers.Home).Methods(http.MethodGet)
	s.HandleFunc("/repositories", controllers.Repositories).Methods(http.MethodGet)
	s.HandleFunc("/create-demoscene", controllers.CreateDemoscene).Methods(http.MethodPost)

	return r
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := encryption.GetLoggedUser(r)
		if err != nil {
			controllers.Login(w, r)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
