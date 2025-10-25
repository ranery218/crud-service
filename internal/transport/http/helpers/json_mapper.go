package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func DecodeJSON[T any](r *http.Request, dst *T) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		switch {
		case errors.Is(err, io.EOF):
			return ErrBodyEmpty

		default:
			var syntaxErr *json.SyntaxError
			if errors.As(err, &syntaxErr) {
				return ErrMalformedJSON
			}

			var typeErr *json.UnmarshalTypeError
			if errors.As(err, &typeErr) {
				return ErrInvalidType
			}
			return err
		}
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}
