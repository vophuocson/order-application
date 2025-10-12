package apiutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	applicationoutbound "user-domain/internal/application/outbound"
	domainoutport "user-domain/internal/domain/outport"

	"github.com/pkg/errors"
)

type JSONResponse interface {
	Success(status int, response interface{}) error
	Failure(err error)
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

func (v *jsonResponse) Failure(e error) {
	v.logger.Info(e.Error(), domainoutport.LogFields{})
	v.w.WriteHeader(http.StatusBadRequest)
}
