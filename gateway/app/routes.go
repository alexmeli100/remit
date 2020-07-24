package app

import "github.com/gorilla/mux"

func (a *App) InitializeRoutes(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/signup", a.createUser()).Methods("POST")
	api.HandleFunc("/signin", a.signIn()).Methods("POST")

	// authentication routes
	s := api.PathPrefix("/auth").Subrouter()
	//s.Use(a.isAuthenticated)
	s.HandleFunc("/user/{id:[0-9]+}", a.getUser()).Methods("GET")
}
