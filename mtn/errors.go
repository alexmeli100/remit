package mtn

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type ErrorHandlerFunc func(res *http.Response) error

const (
	PayeeNotFound               = "PAYEE_NOT_FOUND"
	PayerNotFound               = "PAYER_NOT_FOUND"
	NotAllowed                  = "NOT_ALLOWED"
	NotAllowedTargetEnvironment = "NOT_ALLOWED_TARGET_ENVIRONMENT"
	InvalidCallbackURLHost      = "INVALID_CALLBACK_URL_HOST"
	InvalidCurrency             = "INVALID_CURRENCY"
	ServiceUnavailable          = "SERVICE_UNAVAILABLE"
	InternalProcessing          = "INTERNAL_PROCESSING_ERROR"
	NotEnoughFunds              = "NOT_ENOUGH_FUNDS"
	PayerLimitReached           = "PAYER_LIMIT_REACHED"
	PayeeNotAllowedToReceive    = "PAYEE_NOT_ALLOWED_TO_RECEIVE"
	PaymentNotApproved          = "PAYMENT_NOT_APPROVED"
	ResourceNotFound            = "RESOURCE_NOT_FOUND"
	ApprovalRejected            = "APPROVAL_REJECTED"
	Expired                     = "EXPIRED"
	TransactionCancelled        = "TRANSACTION_CANCELED"
	ResourceAlreadyExists       = "RESOURCE_ALREADY_EXIST"
)

var (
	UnSpecifiedError = errors.New("error not specied")
	ErrorNoBody      = errors.New("no body in response")
	// this error is return when getting the transfer status takes longer than max time
	ErrorPending = errors.New("transfer still pending")
)

type ErrTransferFailed struct {
	message string
}

func (e ErrTransferFailed) Error() string {
	return e.message
}

type ApprovalRejectedError struct {
	message string
}

func (a ApprovalRejectedError) Error() string {
	return a.message
}

type ExpiredError struct {
	message string
}

func (a ExpiredError) Error() string {
	return a.message
}

type InternalProcessingError struct {
	message string
}

func (a InternalProcessingError) Error() string {
	return a.message
}

type InvalidCallbackURLHostError struct {
	message string
}

func (a InvalidCallbackURLHostError) Error() string {
	return a.message
}

type InvalidCurrencyError struct {
	message string
}

func (a InvalidCurrencyError) Error() string {
	return a.message
}

type NotAllowedError struct {
	message string
}

func (a NotAllowedError) Error() string {
	return a.message
}

type NotAllowedTargetEnvironmentError struct {
	message string
}

func (a NotAllowedTargetEnvironmentError) Error() string {
	return a.message
}

type NotEnoughFundsError struct {
	message string
}

func (a NotEnoughFundsError) Error() string {
	return a.message
}

type PayeeNotFoundError struct {
	message string
}

func (a PayeeNotFoundError) Error() string {
	return a.message
}

type PayeeNotAllowedToReceiveError struct {
	message string
}

func (a PayeeNotAllowedToReceiveError) Error() string {
	return a.message
}

type PayerLimitreachedError struct {
	message string
}

func (a PayerLimitreachedError) Error() string {
	return a.message
}

type PayerNotFoundError struct {
	message string
}

func (a PayerNotFoundError) Error() string {
	return a.message
}

type PaymentNotApprovedError struct {
	message string
}

func (a PaymentNotApprovedError) Error() string {
	return a.message
}

type ResourceAlreadyExistsError struct {
	message string
}

func (a ResourceAlreadyExistsError) Error() string {
	return a.message
}

type ResourceNotFoundError struct {
	message string
}

func (a ResourceNotFoundError) Error() string {
	return a.message
}

type ServiceUnavailableError struct {
	message string
}

func (a ServiceUnavailableError) Error() string {
	return a.message
}

type TransactionCancelledError struct {
	message string
}

func (a TransactionCancelledError) Error() string {
	return a.message
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorHandler struct {
	handler ErrorHandlerFunc
	next    ResponseHandler
}

func (e *ErrorHandler) handleResponse(res *http.Response, i interface{}) error {
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return e.next.handleResponse(res, i)
	} else if res.StatusCode == http.StatusUnauthorized {
		return errors.Errorf("status: 401\nmessage: status unauthorized")
	}

	return e.handler(res)
}

func remittanceErrHandler(res *http.Response) error {
	var er *ResponseError
	err := getErrorRes(res, er)

	if err != nil {
		if errors.Is(err, ErrorNoBody) {
			return UnSpecifiedError
		} else {
			return err
		}
	}

	return getError(er, res.StatusCode)
}

func tokenErrHandler(res *http.Response) error {
	var er struct {
		Err string `json:"error"`
	}
	err := getErrorRes(res, &er)

	if err != nil {
		if errors.Is(err, ErrorNoBody) {
			return UnSpecifiedError
		} else {
			return err
		}
	}

	return errors.New(er.Err)
}

func getError(err *ResponseError, status int) error {
	message := fmt.Sprintf("status: %d\nmessage: %s\n", status, err.Message)

	switch err.Code {
	case ApprovalRejected:
		return ApprovalRejectedError{message}
	case Expired:
		return ExpiredError{message}
	case InternalProcessing:
		return InternalProcessingError{message}
	case InvalidCallbackURLHost:
		return InvalidCallbackURLHostError{message}
	case InvalidCurrency:
		return InvalidCurrencyError{message}
	case NotAllowed:
		return NotAllowedError{message}
	case NotAllowedTargetEnvironment:
		return NotAllowedTargetEnvironmentError{message}
	case NotEnoughFunds:
		return NotEnoughFundsError{message}
	case PayeeNotFound:
		return PayeeNotFoundError{message}
	case PayeeNotAllowedToReceive:
		return PayeeNotAllowedToReceiveError{message}
	case PayerLimitReached:
		return PayerLimitreachedError{message}
	case PayerNotFound:
		return PayeeNotFoundError{message}
	case PaymentNotApproved:
		return PaymentNotApprovedError{message}
	case ResourceAlreadyExists:
		return ResourceAlreadyExistsError{message}
	case ResourceNotFound:
		return ResourceNotFoundError{message}
	case ServiceUnavailable:
		return ServiceUnavailableError{message}
	case TransactionCancelled:
		return TransactionCancelledError{message}

	default:
		return errors.New("unknown error")
	}
}

func getErrorRes(res *http.Response, i interface{}) error {
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return errors.Wrap(err, "error reading response")
	}

	if len(body) > 0 {

		if err = json.Unmarshal(body, i); err != nil {
			return errors.Wrap(err, "json unmarshal error")
		}

		return nil
	}

	return ErrorNoBody
}
