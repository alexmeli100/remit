package app

import "github.com/gorilla/mux"

func (a *App) InitializeRoutes(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()
	web := api.PathPrefix("/web").Subrouter()
	mobile := api.PathPrefix("/mobile").Subrouter()
	api.HandleFunc("/signup", a.createUser()).Methods("POST")
	web.HandleFunc("/signin", a.signInWeb()).Methods("POST")

	mobileAuth := mobile.PathPrefix("/auth").Subrouter()
	mobileAuth.Use(a.isAuthenticatedMobile)

	webAuth := web.PathPrefix("/auth").Subrouter()
	webAuth.Use(a.isAuthenticatedWeb)
	clients := []*mux.Router{mobileAuth, webAuth}

	// authentication routes
	for _, r := range clients {
		r.HandleFunc("/user/uid/{id:[0-9]+}", a.getUserByID()).Methods("GET")
		r.HandleFunc("/user/id/{id:[a-zA-Z0-9]+}", a.getUserByUUID()).Methods("GET")
	}
}
