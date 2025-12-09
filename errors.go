package action

import (
	"errors"
	"log/slog"
	"net/http"
)

type Error struct {
	code     int
	internal error
}

func NewError(code int) Error {
	return Error{code: code}
}

func (err Error) Code() int { return err.code }

func (err Error) Error() string { return http.StatusText(err.code) }

func (err Error) Unwrap() error { return err.internal }

func (err Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("code", err.code),
		slog.String("status", http.StatusText(err.code)),
		slog.Any("internal", err.internal),
	)
}

func (err Error) Render(w http.ResponseWriter, r *http.Request) error {
	http.Error(w, err.Error(), err.Code())
	return nil
}

func (err Error) WithInternal(internal error) Error {
	err.internal = internal
	return err
}

var (
	ErrBadRequest                    = NewError(http.StatusBadRequest)                    // RFC 9110, 15.5.1
	ErrUnauthorized                  = NewError(http.StatusUnauthorized)                  // RFC 9110, 15.5.2
	ErrPaymentRequired               = NewError(http.StatusPaymentRequired)               // RFC 9110, 15.5.3
	ErrForbidden                     = NewError(http.StatusForbidden)                     // RFC 9110, 15.5.4
	ErrNotFound                      = NewError(http.StatusNotFound)                      // RFC 9110, 15.5.5
	ErrMethodNotAllowed              = NewError(http.StatusMethodNotAllowed)              // RFC 9110, 15.5.6
	ErrNotAcceptable                 = NewError(http.StatusNotAcceptable)                 // RFC 9110, 15.5.7
	ErrProxyAuthRequired             = NewError(http.StatusProxyAuthRequired)             // RFC 9110, 15.5.8
	ErrRequestTimeout                = NewError(http.StatusRequestTimeout)                // RFC 9110, 15.5.9
	ErrConflict                      = NewError(http.StatusConflict)                      // RFC 9110, 15.5.10
	ErrGone                          = NewError(http.StatusGone)                          // RFC 9110, 15.5.11
	ErrLengthRequired                = NewError(http.StatusLengthRequired)                // RFC 9110, 15.5.12
	ErrPreconditionFailed            = NewError(http.StatusPreconditionFailed)            // RFC 9110, 15.5.13
	ErrRequestEntityTooLarge         = NewError(http.StatusRequestEntityTooLarge)         // RFC 9110, 15.5.14
	ErrRequestURITooLong             = NewError(http.StatusRequestURITooLong)             // RFC 9110, 15.5.15
	ErrUnsupportedMediaType          = NewError(http.StatusUnsupportedMediaType)          // RFC 9110, 15.5.16
	ErrRequestedRangeNotSatisfiable  = NewError(http.StatusRequestedRangeNotSatisfiable)  // RFC 9110, 15.5.17
	ErrExpectationFailed             = NewError(http.StatusExpectationFailed)             // RFC 9110, 15.5.18
	ErrTeapot                        = NewError(http.StatusTeapot)                        // RFC 9110, 15.5.19 (Unused)
	ErrMisdirectedRequest            = NewError(http.StatusMisdirectedRequest)            // RFC 9110, 15.5.20
	ErrUnprocessableEntity           = NewError(http.StatusUnprocessableEntity)           // RFC 9110, 15.5.21
	ErrLocked                        = NewError(http.StatusLocked)                        // RFC 4918, 11.3
	ErrFailedDependency              = NewError(http.StatusFailedDependency)              // RFC 4918, 11.4
	ErrTooEarly                      = NewError(http.StatusTooEarly)                      // RFC 8470, 5.2.
	ErrUpgradeRequired               = NewError(http.StatusUpgradeRequired)               // RFC 9110, 15.5.22
	ErrPreconditionRequired          = NewError(http.StatusPreconditionRequired)          // RFC 6585, 3
	ErrTooManyRequests               = NewError(http.StatusTooManyRequests)               // RFC 6585, 4
	ErrRequestHeaderFieldsTooLarge   = NewError(http.StatusRequestHeaderFieldsTooLarge)   // RFC 6585, 5
	ErrUnavailableForLegalReasons    = NewError(http.StatusUnavailableForLegalReasons)    // RFC 7725, 3
	ErrInternalServerError           = NewError(http.StatusInternalServerError)           // RFC 9110, 15.6.1
	ErrNotImplemented                = NewError(http.StatusNotImplemented)                // RFC 9110, 15.6.2
	ErrBadGateway                    = NewError(http.StatusBadGateway)                    // RFC 9110, 15.6.3
	ErrServiceUnavailable            = NewError(http.StatusServiceUnavailable)            // RFC 9110, 15.6.4
	ErrGatewayTimeout                = NewError(http.StatusGatewayTimeout)                // RFC 9110, 15.6.5
	ErrHTTPVersionNotSupported       = NewError(http.StatusHTTPVersionNotSupported)       // RFC 9110, 15.6.6
	ErrVariantAlsoNegotiates         = NewError(http.StatusVariantAlsoNegotiates)         // RFC 2295, 8.1
	ErrInsufficientStorage           = NewError(http.StatusInsufficientStorage)           // RFC 4918, 11.5
	ErrLoopDetected                  = NewError(http.StatusLoopDetected)                  // RFC 5842, 7.2
	ErrNotExtended                   = NewError(http.StatusNotExtended)                   // RFC 2774, 7
	ErrNetworkAuthenticationRequired = NewError(http.StatusNetworkAuthenticationRequired) // RFC 6585, 6
)

var HandleError = DefaultErrorHandler

func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("request failed", "err", err)

	var actionErr Error
	if !errors.As(err, &actionErr) {
		actionErr = ErrInternalServerError
	}

	_ = actionErr.Render(w, r)
}
