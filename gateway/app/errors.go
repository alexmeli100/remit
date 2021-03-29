package app

import (
	"net/http"
)

func (a *App) serverError(w http.ResponseWriter, err error) {
	a.respondWithError(w, http.StatusInternalServerError, err)
}

func (a *App) badRequest(w http.ResponseWriter, err error) {
	a.respondWithError(w, http.StatusBadRequest, err)
}

func (a *App) unauthorized(w http.ResponseWriter, err error) {
	a.respondWithError(w, http.StatusUnauthorized, err)
}

func (a *App) notFound(w http.ResponseWriter, err error) {
	a.respondWithError(w, http.StatusNotFound, err)
}
