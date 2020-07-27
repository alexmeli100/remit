package app

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/alexmeli100/remit/events"
	notificator "github.com/alexmeli100/remit/notificator/pkg/service"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	user "github.com/alexmeli100/remit/users/pkg/service"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io"
	log1 "log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type App struct {
	Server       *http.Server
	Events       events.EventManager
	UsersService user.UsersService
	Notificator  notificator.NotificatorService
	FireApp      *firebase.App
	Logger       log.Logger
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

func (a *App) isAuthenticatedMobile(next http.Handler) http.Handler {
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
		}

		token, err := client.VerifyIDToken(r.Context(), idToken)
		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// create user in firebase then add the user to our database.
// an account activation email is sent afterwards.
func (a *App) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := decodeBody(r.Body)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		name := fmt.Sprintf("%s %s", u.FirstName, u.LastName)
		if err = a.UsersService.Create(r.Context(), u); err != nil {
			a.serverError(w, err)
			return
		}
		respondWithJson(w, http.StatusCreated, map[string]string{"status": "user created"})
		a.Logger.Log("method", "createUser", "name", name)

		if err := a.Events.OnUserCreated(r.Context(), u); err != nil {
			a.Logger.Log("method", "createUser", "err", err)
		}
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

		if err != nil {
			a.serverError(w, err)
			return
		}

		respondWithJson(w, http.StatusOK, u)
	}
}

func (a *App) getUserByUUID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		u, err := a.UsersService.GetUserByUUID(r.Context(), id)

		if err != nil {
			a.serverError(w, err)
			return
		}

		respondWithJson(w, http.StatusOK, u)
	}
}

func (a *App) signIn() http.HandlerFunc {
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

		if time.Now().Unix()-decoded.Claims["auth_time"].(int64) > 5*60 {
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
		respondWithJson(w, http.StatusOK, response)
	}
}

func (a *App) OnUserCreated(ctx context.Context, u *pb.User) error {
	client, err := a.FireApp.Auth(ctx)

	if err != nil {
		return err
	}

	url, err := client.EmailVerificationLink(ctx, u.Email)

	if err != nil {
		return errors.Wrap(err, "error getting email confirmation link")
	}

	if err = a.Notificator.SendConfirmEmail(ctx, u.FirstName, u.Email, url); err != nil {
		return errors.Wrap(err, "error sending confirmation email")
	}

	if err = a.Notificator.SendWelcomeEmail(ctx, u.FirstName, u.Email); err != nil {
		return errors.Wrap(err, "error sending welcome email")
	}

	return nil
}

func (a *App) OnPasswordReset(ctx context.Context, u *pb.User) error {
	client, err := a.FireApp.Auth(ctx)

	if err != nil {
		return err
	}

	url, err := client.PasswordResetLink(ctx, u.Email)

	if err != nil {
		return errors.Wrap(err, "error getting password reset link")
	}

	if err = a.Notificator.SendPasswordResetEmail(ctx, u.Email, url); err != nil {
		return errors.Wrap(err, "error sending password reset link")
	}

	return nil
}

// decode body of a request containing user data.
// used in create user method
func decodeBody(body io.Reader) (*pb.User, error) {
	var u pb.User
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&u); err != nil {
		return nil, errors.Wrap(err, "error decoding request")
	}

	return &u, nil
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

func respondWithError(w http.ResponseWriter, code int, err error) {
	respondWithJson(w, code, map[string]string{"error": err.Error()})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log1.Println(err)
	}
}
