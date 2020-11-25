package app

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/alexmeli100/remit/events"
	notificator "github.com/alexmeli100/remit/notificator/pkg/service"
	paymentpb "github.com/alexmeli100/remit/payment/pkg/grpc/pb"
	payment "github.com/alexmeli100/remit/payment/pkg/service"
	transferpb "github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
	transfer "github.com/alexmeli100/remit/transfer/pkg/service"
	userpb "github.com/alexmeli100/remit/users/pkg/grpc/pb"
	user "github.com/alexmeli100/remit/users/pkg/service"
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v71"
	stripeclient "github.com/stripe/stripe-go/v71/client"
	"github.com/stripe/stripe-go/v71/webhook"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type App struct {
	RedisClient     *redis.Client
	StripeClient    *stripeclient.API
	Server          *http.Server
	Events          events.EventManager
	UsersService    user.UsersService
	PaymentService  payment.PaymentService
	TransferService transfer.TransferService
	Notificator     notificator.NotificatorService
	FireApp         *firebase.App
	Logger          log.Logger
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

func (a *App) checkWebHookSignature(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		maxBodyBytes := int64(65536)
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		payload, err := ioutil.ReadAll(r.Body)

		if err != nil {
			err = errors.Wrap(err, "error reading request body")
			a.Logger.Log("error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		endpointSecret := os.Getenv("STRIPE_ENDPOINT_SECRET")
		signature := r.Header.Get("Stripe-Signature")
		event, err := webhook.ConstructEvent(payload, signature, endpointSecret)

		if err != nil {
			err = errors.Wrap(err, "failed to verify webhook signature")
			a.Logger.Log("error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "event", &event)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// create user in firebase then add the user to our database.
// an account activation email is sent afterwards.
func (a *App) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &userpb.User{}
		err := decodeRequestBody(r.Body, u)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		cu, err := a.UsersService.Create(r.Context(), u)

		if err != nil {
			a.serverError(w, err)
			return
		}

		a.respondWithJson(w, http.StatusCreated, cu)

		if err := a.Events.OnUserCreated(r.Context(), u); err != nil {
			a.Logger.Log("method", "events.OnUserCreated", "err", err)
		}
	}
}

func (a *App) createContact() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &userpb.Contact{}
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

		if err = a.UsersService.DeleteContact(r.Context(), &userpb.Contact{Id: int64(id)}); err != nil {
			a.serverError(w, err)
			return
		}

		a.respondWithJson(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

func (a *App) createTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tr := &paymentpb.Transaction{}

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

func (a *App) getCustomerID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]

		cid, err := a.PaymentService.GetCustomerID(r.Context(), uid)

		if err != nil {
			a.serverError(w, err)
			return
		}

		res := map[string]string{"customerID": cid}
		a.respondWithJson(w, http.StatusOK, res)
	}
}

func (a *App) setUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &userpb.User{}

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
		u := &userpb.User{}

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
		c := &userpb.Contact{}

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

func (a *App) stripeWebHook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event := r.Context().Value("event").(*stripe.Event)
		err := a.checkIfPaymentEventProcessed(r.Context(), event)

		if errors.Is(err, ErrorEventProcessed) {
			w.WriteHeader(http.StatusOK)
			return
		} else if err != nil {
			a.Logger.Log("method", "stripeWebHook", "error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		switch event.Type {
		case "payment_intent.succeeded":
			if err := a.handlePaymentSucceded(r.Context(), w, event); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
		default:
			a.Logger.Log("error", fmt.Sprintf("unexpected event type: %s\n", event.Type))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (a *App) handlePaymentSucceded(ctx context.Context, w http.ResponseWriter, e *stripe.Event) error {
	var intent stripe.PaymentIntent

	if err := json.Unmarshal(e.Data.Raw, &intent); err != nil {
		err = errors.Wrap(err, "error parsing webhook json")
		a.Logger.Log("error", err)
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	if err := a.Events.OnPaymentSucceded(ctx, intent.ID); err != nil {
		a.Logger.Log("error", err)
		return err
	}

	return nil
}

func (a *App) checkIfPaymentEventProcessed(ctx context.Context, e *stripe.Event) error {
	var cachedEvent stripe.Event
	err := a.getKey(ctx, e.ID, cachedEvent)

	if err == redis.Nil {
		if err := a.setKey(ctx, e.ID, e, time.Hour*24); err != nil {
			return err
		}

		return nil
	} else if err != nil {
		return err
	}

	return ErrorEventProcessed
}

func (a *App) checkGetUserError(err error, w http.ResponseWriter, u *userpb.User) {
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
		a.Logger.Log("authTime", authTime)

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
		u := &userpb.User{}
		err := decodeRequestBody(r.Body, u)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		err = a.Events.OnPasswordReset(r.Context(), u)
		a.Logger.Log("method", "resetPassword", "error", err)
		a.respondWithJson(w, http.StatusOK, nil)
	}
}

func (a *App) transferMoney() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := &transferpb.TransferRequest{}
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
		if err = a.Events.OnTransferSucceded(r.Context(), res); err != nil {
			a.Logger.Log("method", "transferMoney", "error", err)
		}
	}
}

func (a *App) saveUserCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &userpb.User{}
		err := decodeRequestBody(r.Body, u)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		s, err := a.PaymentService.SaveCard(r.Context(), u.Uuid)

		if err != nil {
			a.serverError(w, err)
			return
		}

		a.respondWithJson(w, http.StatusOK, map[string]string{"secret": s})
	}
}

func (a *App) getPaymentSecret() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &paymentpb.PaymentRequest{}
		err := decodeRequestBody(r.Body, p)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		s, err := a.PaymentService.GetPaymentIntentSecret(r.Context(), p)

		if err != nil {
			a.serverError(w, err)
			return
		}

		a.respondWithJson(w, http.StatusOK, map[string]string{"secret": s})
	}
}

func (a *App) capturePayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var c struct {
			Pi     string  `json:"pi"`
			Amount float64 `json:"amount"`
		}
		err := decodeRequestBody(r.Body, &c)

		if err != nil {
			a.badRequest(w, err)
			return
		}

		s, err := a.PaymentService.CapturePayment(r.Context(), c.Pi, c.Amount)

		if err != nil {
			a.serverError(w, err)
			return
		}

		a.respondWithJson(w, http.StatusOK, map[string]string{"secret": s})
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
	a.Logger.Log("code", code, "error", err)
	a.respondWithJson(w, code, map[string]string{"error": err.Error()})
}

func (a *App) respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			a.Logger.Log("error", err)
		}
	}

}

func (a *App) getKey(ctx context.Context, key string, src interface{}) error {
	val, err := a.RedisClient.Get(ctx, key).Result()

	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(val), src); err != nil {
		return err
	}

	return nil
}

func (a *App) setKey(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	v, err := json.Marshal(value)

	if err != nil {
		return err
	}

	err = a.RedisClient.Set(ctx, key, v, expiration).Err()

	if err != nil {
		return err
	}

	return nil
}
