package rho

import (
	"errors"
	"log/slog"
	"net/http"
)

type HttpError struct {
	code     int
	internal error
}

func NewHttpError(code int) HttpError {
	return HttpError{code: code}
}

func (err HttpError) Code() int { return err.code }

func (err HttpError) Error() string { return http.StatusText(err.code) }

func (err HttpError) Unwrap() error { return err.internal }

func (err HttpError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("code", err.code),
		slog.String("status", http.StatusText(err.code)),
		slog.Any("internal", err.internal),
	)
}

func (err HttpError) Render(w http.ResponseWriter, r *http.Request) error {
	http.Error(w, err.Error(), err.Code())
	return nil
}

func (err HttpError) WithInternal(internal error) HttpError {
	err.internal = internal
	return err
}

var (
	ErrBadRequest                    = NewHttpError(http.StatusBadRequest)                    // RFC 9110, 15.5.1
	ErrUnauthorized                  = NewHttpError(http.StatusUnauthorized)                  // RFC 9110, 15.5.2
	ErrPaymentRequired               = NewHttpError(http.StatusPaymentRequired)               // RFC 9110, 15.5.3
	ErrForbidden                     = NewHttpError(http.StatusForbidden)                     // RFC 9110, 15.5.4
	ErrNotFound                      = NewHttpError(http.StatusNotFound)                      // RFC 9110, 15.5.5
	ErrMethodNotAllowed              = NewHttpError(http.StatusMethodNotAllowed)              // RFC 9110, 15.5.6
	ErrNotAcceptable                 = NewHttpError(http.StatusNotAcceptable)                 // RFC 9110, 15.5.7
	ErrProxyAuthRequired             = NewHttpError(http.StatusProxyAuthRequired)             // RFC 9110, 15.5.8
	ErrRequestTimeout                = NewHttpError(http.StatusRequestTimeout)                // RFC 9110, 15.5.9
	ErrConflict                      = NewHttpError(http.StatusConflict)                      // RFC 9110, 15.5.10
	ErrGone                          = NewHttpError(http.StatusGone)                          // RFC 9110, 15.5.11
	ErrLengthRequired                = NewHttpError(http.StatusLengthRequired)                // RFC 9110, 15.5.12
	ErrPreconditionFailed            = NewHttpError(http.StatusPreconditionFailed)            // RFC 9110, 15.5.13
	ErrRequestEntityTooLarge         = NewHttpError(http.StatusRequestEntityTooLarge)         // RFC 9110, 15.5.14
	ErrRequestURITooLong             = NewHttpError(http.StatusRequestURITooLong)             // RFC 9110, 15.5.15
	ErrUnsupportedMediaType          = NewHttpError(http.StatusUnsupportedMediaType)          // RFC 9110, 15.5.16
	ErrRequestedRangeNotSatisfiable  = NewHttpError(http.StatusRequestedRangeNotSatisfiable)  // RFC 9110, 15.5.17
	ErrExpectationFailed             = NewHttpError(http.StatusExpectationFailed)             // RFC 9110, 15.5.18
	ErrTeapot                        = NewHttpError(http.StatusTeapot)                        // RFC 9110, 15.5.19 (Unused)
	ErrMisdirectedRequest            = NewHttpError(http.StatusMisdirectedRequest)            // RFC 9110, 15.5.20
	ErrUnprocessableEntity           = NewHttpError(http.StatusUnprocessableEntity)           // RFC 9110, 15.5.21
	ErrLocked                        = NewHttpError(http.StatusLocked)                        // RFC 4918, 11.3
	ErrFailedDependency              = NewHttpError(http.StatusFailedDependency)              // RFC 4918, 11.4
	ErrTooEarly                      = NewHttpError(http.StatusTooEarly)                      // RFC 8470, 5.2.
	ErrUpgradeRequired               = NewHttpError(http.StatusUpgradeRequired)               // RFC 9110, 15.5.22
	ErrPreconditionRequired          = NewHttpError(http.StatusPreconditionRequired)          // RFC 6585, 3
	ErrTooManyRequests               = NewHttpError(http.StatusTooManyRequests)               // RFC 6585, 4
	ErrRequestHeaderFieldsTooLarge   = NewHttpError(http.StatusRequestHeaderFieldsTooLarge)   // RFC 6585, 5
	ErrUnavailableForLegalReasons    = NewHttpError(http.StatusUnavailableForLegalReasons)    // RFC 7725, 3
	ErrInternalServerError           = NewHttpError(http.StatusInternalServerError)           // RFC 9110, 15.6.1
	ErrNotImplemented                = NewHttpError(http.StatusNotImplemented)                // RFC 9110, 15.6.2
	ErrBadGateway                    = NewHttpError(http.StatusBadGateway)                    // RFC 9110, 15.6.3
	ErrServiceUnavailable            = NewHttpError(http.StatusServiceUnavailable)            // RFC 9110, 15.6.4
	ErrGatewayTimeout                = NewHttpError(http.StatusGatewayTimeout)                // RFC 9110, 15.6.5
	ErrHTTPVersionNotSupported       = NewHttpError(http.StatusHTTPVersionNotSupported)       // RFC 9110, 15.6.6
	ErrVariantAlsoNegotiates         = NewHttpError(http.StatusVariantAlsoNegotiates)         // RFC 2295, 8.1
	ErrInsufficientStorage           = NewHttpError(http.StatusInsufficientStorage)           // RFC 4918, 11.5
	ErrLoopDetected                  = NewHttpError(http.StatusLoopDetected)                  // RFC 5842, 7.2
	ErrNotExtended                   = NewHttpError(http.StatusNotExtended)                   // RFC 2774, 7
	ErrNetworkAuthenticationRequired = NewHttpError(http.StatusNetworkAuthenticationRequired) // RFC 6585, 6
)

var HandleError = DefaultErrorHandler

func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("request failed", "err", err)

	var herr HttpError
	if !errors.As(err, &herr) {
		herr = ErrInternalServerError
	}

	_ = herr.Render(w, r)
}
