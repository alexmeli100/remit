package app

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	notificator "github.com/alexmeli100/remit/notificator/pkg/service"
	payment "github.com/alexmeli100/remit/payment/pkg/service"
	transfer "github.com/alexmeli100/remit/transfer/pkg/service"
	user "github.com/alexmeli100/remit/users/pkg/service"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type App struct {
	Server          *http.Server
	UsersService    user.UsersService
	PaymentService  payment.PaymentService
	TransferService transfer.TransferService
	Notificator     notificator.NotificatorService
	FireApp         *firebase.App
	Logger          *zap.Logger
}

func (a *App) isAuthenticatedWeb(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")

		if err != nil {
			err = errors.New("session cookie unavailable")
			a.unauthorized(w, err)
			return
		}

		token, err := getTokenFromSession(r.Context(), a.FireApp, cookie.Value)

		if err != nil {
			a.serverError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *App) isAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client, err := a.FireApp.Auth(r.Context())

		if err != nil {
			a.serverError(w, err)
			return
		}

		bearer := "Bearer "
		authHeader := r.Header.Get("Authorization")
		var idToken string

		if strings.HasPrefix(authHeader, bearer) {
			idToken = authHeader[len(bearer):]
		} else {
			a.badRequest(w, errors.New("Invalid bearer token"))
			return
		}

		token, err := client.VerifyIDToken(r.Context(), idToken)

		if err != nil {
			a.badRequest(w, errors.New("invalid firebase token"))
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// create user in firebase then add the user to our database.
// an account activation email is sent afterwards.
func (a *App) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &user.User{}
		err := decodeRequestBody(r.Body, u)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		cu, err := a.UsersService.CreateUser(r.Context(), u)

		if err != nil {
			a.serverError(w, err)
			return
		}

		a.respondWithJson(w, http.StatusCreated, cu)

		if err := a.sendWelcomeEmail(r.Context(), cu); err != nil {
			a.Logger.Info("Error sending welcome emails", zap.Error(err))
		}
	}
}

func (a *App) sendWelcomeEmail(ctx context.Context, u *user.User) error {
	client, err := a.FireApp.Auth(ctx)

	if err != nil {
		return err
	}

	url, err := client.EmailVerificationLink(ctx, u.Email)

	if err != nil {
		return errors.Wrap(err, "Error getting email verification link")
	}

	if err = a.Notificator.SendConfirmEmail(ctx, u.FirstName, u.Email, url); err != nil {
		return errors.Wrap(err, "error sending confirmation email")
	}

	if err = a.Notificator.SendWelcomeEmail(ctx, u.FirstName, u.Email); err != nil {
		return errors.Wrap(err, "error sending welcome email")
	}

	return nil
}

func (a *App) createContact() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &user.Contact{}
		err := decodeRequestBody(r.Body, c)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		uc, err := a.UsersService.CreateContact(r.Context(), c)

		if err != nil {
			a.serverError(w, err)
			return
		}

		a.respondWithJson(w, http.StatusCreated, uc)
	}
}

func (a *App) getContacts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		p := vars["id"]

		id, err := strconv.Atoi(p)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		contacts, err := a.UsersService.GetContacts(r.Context(), int64(id))

		if err != nil {
			a.serverError(w, err)
			return
		}

		a.respondWithJson(w, http.StatusCreated, contacts)
	}
}

func (a *App) deleteContact() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		p := vars["id"]

		id, err := strconv.Atoi(p)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		if err = a.UsersService.DeleteContact(r.Context(), &user.Contact{Id: int64(id)}); err != nil {
			a.serverError(w, err)
			return
		}

		a.respondWithJson(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

func (a *App) createTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tr := &payment.Transaction{}

		if err := decodeRequestBody(r.Body, tr); err != nil {
			a.badRequest(w, err)
			return
		}

		tr, err := a.PaymentService.CreateTransaction(r.Context(), tr)

		if err != nil {
			a.serverError(w, err)
			return
		}

		a.respondWithJson(w, http.StatusCreated, tr)
	}
}

func (a *App) getTransactions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]

		trs, err := a.PaymentService.GetTransactions(r.Context(), uid)

		if err != nil {
			a.serverError(w, err)
			return
		}

		a.respondWithJson(w, http.StatusOK, trs)
	}
}

func (a *App) setUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &user.User{}

		if err := decodeRequestBody(r.Body, u); err != nil {
			a.badRequest(w, err)
			return
		}

		u, err := a.UsersService.SetUserProfile(r.Context(), u)

		if err != nil {
			a.serverError(w, err)
			return
		}

		a.respondWithJson(w, http.StatusOK, u)
	}
}

func (a *App) updateUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &user.User{}

		if err := decodeRequestBody(r.Body, u); err != nil {
			a.badRequest(w, err)
			return
		}

		u, err := a.UsersService.UpdateUserProfile(r.Context(), u)

		if err != nil {
			a.serverError(w, err)
			return
		}

		a.respondWithJson(w, http.StatusOK, u)
	}
}

func (a *App) updateContact() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &user.Contact{}

		err := decodeRequestBody(r.Body, c)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		uc, err := a.UsersService.UpdateContact(r.Context(), c)

		if err != nil {
			a.serverError(w, err)
			return
		}

		a.respondWithJson(w, http.StatusOK, uc)
	}
}

func (a *App) getUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		p := vars["id"]
		id, err := strconv.Atoi(p)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		u, err := a.UsersService.GetUserByID(r.Context(), int64(id))

		a.checkGetUserError(err, w, u)
	}
}

func (a *App) getUserByUUID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]

		u, err := a.UsersService.GetUserByUUID(r.Context(), uid)

		a.checkGetUserError(err, w, u)
	}
}

func (a *App) checkGetUserError(err error, w http.ResponseWriter, u *user.User) {
	if errors.Is(err, user.ErrUserNotFound) {
		a.notFound(w, err)
		return
	} else if err != nil {
		a.serverError(w, err)
		return
	}

	a.respondWithJson(w, http.StatusOK, u)
}

func (a *App) signInWeb() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idToken, err := getIdToken(r)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		expiresIn := time.Hour * 24

		client, err := a.FireApp.Auth(r.Context())

		if err != nil {
			a.serverError(w, err)
			return
		}

		decoded, err := client.VerifyIDToken(r.Context(), idToken)

		if err != nil {
			err = errors.Wrap(err, "invalid ID token")
			a.unauthorized(w, err)
			return
		}

		authTime := time.Now().Unix() - int64(decoded.Claims["auth_time"].(float64))
		a.Logger.Info("authTime", zap.Int64("auth time difference", authTime))

		if authTime > 5*60 {
			err = errors.New("recent sign-in required")
			a.unauthorized(w, err)
			return
		}

		cookie, err := client.SessionCookie(r.Context(), idToken, expiresIn)

		if err != nil {
			err = errors.Wrap(err, "failed to create session cookie")
			a.serverError(w, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    cookie,
			MaxAge:   int(expiresIn),
			HttpOnly: true,
			Secure:   true,
		})

		response := map[string]string{"status": "success"}
		a.respondWithJson(w, http.StatusOK, response)
	}
}

func (a *App) resetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &user.User{}
		err := decodeRequestBody(r.Body, u)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		client, err := a.FireApp.Auth(r.Context())

		if err != nil {
			a.Logger.Error("Error getting firebase client", zap.Error(err))
		}

		url, err := client.PasswordResetLink(r.Context(), u.Email)

		if err = a.Notificator.SendPasswordResetEmail(r.Context(), u.Email, url); err != nil {
			a.Logger.Error("Error sending password reset link", zap.Error(err))
		}

		a.respondWithJson(w, http.StatusOK, nil)
	}
}

func (a *App) transferMoney() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := &transfer.TransferRequest{}
		err := decodeRequestBody(r.Body, t)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		res := a.TransferService.Transfer(r.Context(), t)

		if res.Status == "Failed" {
			a.serverError(w, errors.New(res.FailReason))
			return
		}

		a.respondWithJson(w, http.StatusOK, map[string]string{"ok": "transfer successful"})
	}
}

// decode body of a request containing json data.
func decodeRequestBody(body io.Reader, dst interface{}) error {
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&dst); err != nil {
		return errors.Wrap(err, "error decoding request")
	}

	return nil
}

// get id token from request body
func getIdToken(r *http.Request) (string, error) {
	var token struct {
		IdToken string `json:"idToken"`
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&token); err != nil {
		return "", errors.Wrap(err, "error getting idToken")
	}

	return token.IdToken, nil
}

// verify session is valid and get token
func getTokenFromSession(ctx context.Context, app *firebase.App, idToken string) (*auth.Token, error) {
	client, err := app.Auth(ctx)

	if err != nil {
		return nil, err
	}

	token, err := client.VerifySessionCookieAndCheckRevoked(ctx, idToken)

	if err != nil {
		return nil, err
	}
	return token, nil
}

func (a *App) respondWithError(w http.ResponseWriter, code int, err error) {
	a.Logger.Info("json error response", zap.Int("code", code), zap.Error(err))
	a.respondWithJson(w, code, map[string]string{"error": err.Error()})
}

func (a *App) respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			a.Logger.Error("error sending json response", zap.Error(err))
		}
	}

}
