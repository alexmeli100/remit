package app

import "github.com/gorilla/mux"

func (a *App) InitializeRoutes(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()
	web := api.PathPrefix("/web").Subrouter()
	mobile := api.PathPrefix("/mobile").Subrouter()
	api.PathPrefix("/signup").Handler(a.isAuthenticated(a.createUser())).Methods("POST")
	api.PathPrefix("/reset-password").Handler(a.resetPassword()).Methods("POST")
	//web.HandleFunc("/signin", a.signInWeb()).Methods("POST")

	mobileAuth := mobile.PathPrefix("/auth").Subrouter()
	mobileAuth.Use(a.isAuthenticated)

	webAuth := web.PathPrefix("/auth").Subrouter()
	webAuth.Use(a.isAuthenticated)
	clients := []*mux.Router{mobileAuth, webAuth}

	// common authentication routes
	for _, r := range clients {
		r.HandleFunc("/user/uid/{uid:[a-zA-Z0-9]+}", a.getUserByID()).Methods("GET")
		r.HandleFunc("/user/id/{id:[0-9]+}", a.getUserByUUID()).Methods("GET")
		r.HandleFunc("/user/profile", a.setUserProfile()).Methods("POST")
		r.HandleFunc("/user/profile", a.updateUserProfile()).Methods("PUT")
		r.HandleFunc("/contact", a.createContact()).Methods("POST")
		r.HandleFunc("/contact", a.updateContact()).Methods("PUT")
		r.HandleFunc("/contact/{id:[0-9]+}", a.deleteContact()).Methods("DELETE")
		r.HandleFunc("/contacts/{id:[0-9]+}", a.getContacts()).Methods("GET")
		r.HandleFunc("/transfer", a.transferMoney()).Methods("POST")
		r.HandleFunc("/get-customer-id/uid/{uid:[a-zA-Z0-9]+}", a.getCustomerID()).Methods("GET")
		r.HandleFunc("/create-transaction", a.createTransaction()).Methods("POST")
		r.HandleFunc("/get-transactions/{uid:[a-zA-Z0-9]+}", a.getTransactions()).Methods("GET")
	}
}
