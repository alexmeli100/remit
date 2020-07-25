package app

import "github.com/gorilla/mux"

func (a *App) InitializeRoutes(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/signup", a.createUser()).Methods("POST")
	api.HandleFunc("/signin", a.signIn()).Methods("POST")

	// authentication routes
	s := api.PathPrefix("/auth").Subrouter()
	s.Use(a.isAuthenticated)
	s.HandleFunc("/user/uid/{id:[0-9]+}", a.getUserByID()).Methods("GET")
	s.HandleFunc("/user/id/{id:[a-zA-Z0-9]+}", a.getUserByUUID()).Methods("GET")
}
