package apiutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	applicationerror "user-domain/internal/application/error"
	applicationoutbound "user-domain/internal/application/outbound"
	domainerror "user-domain/internal/domain/error"

	"github.com/pkg/errors"
)

type JSONResponse interface {
	Success(status int, response interface{}) error
	Failure(err error)
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type jsonResponse struct {
	request *http.Request
	w       http.ResponseWriter
	logger  applicationoutbound.Logger
}

func NewJSONResponse(w http.ResponseWriter, r *http.Request, log applicationoutbound.Logger) JSONResponse {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	return &jsonResponse{
		request: r,
		w:       w,
		logger:  log,
	}
}

func (v *jsonResponse) Success(status int, response interface{}) error {
	statusCode, body, err := responseNil(status, response)
	if err != nil {
		v.logger.WithContext(v.request.Context()).Error("apiutil: failed to marshal response: %w", err)
		return fmt.Errorf("apiutil: failed to marshal response: %w", err)
	}

	v.w.WriteHeader(statusCode)
	if body != nil {
		if _, err := v.w.Write(body); err != nil {
			v.logger.WithContext(v.request.Context()).Error("apiutil: failed to write body: %w", err)
			return fmt.Errorf("apiutil: failed to write body: %w", err)
		}
	}

	return nil
}

func responseNil(status int, response interface{}) (int, []byte, error) {
	body, err := json.Marshal(response)
	if err != nil {
		return status, nil, errors.Wrap(err, "apiutil: failed to marshal response")
	}

	if string(body) != "null" {
		return status, body, nil
	}

	if string(body) == "null" && status == http.StatusOK {
		return http.StatusNoContent, nil, nil
	}

	return status, nil, nil
}

func (v *jsonResponse) Failure(err error) {
	var code int
	switch {
	case errors.Is(err, domainerror.ErrCodeNotFound):
		code = http.StatusNotFound
	case errors.Is(err, domainerror.ErrCodeConflict):
		code = http.StatusConflict
	case errors.Is(err, domainerror.ErrCodeInvalidInput),
		errors.Is(err, applicationerror.ErrDecode):
		code = http.StatusBadRequest
	case errors.Is(err, domainerror.ErrCodeForbidden):
		code = http.StatusForbidden
	default:
		code = http.StatusInternalServerError
	}

	v.w.WriteHeader(code)
	if code == http.StatusInternalServerError {
		v.logger.WithContext(v.request.Context()).Error("error: %s", err)
	} else {
		v.logger.WithContext(v.request.Context()).Warn("error: %s", err)
	}

	_ = json.NewEncoder(v.w).Encode(ErrorResponse{
		Code:    code,
		Message: err.Error(),
	})
}
