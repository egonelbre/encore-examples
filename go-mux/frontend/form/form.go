package form

import (
	"net/http"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// Decode decodes the form data into a struct of type T.
func Decode[T any](req *http.Request) (T, error) {
	if err := req.ParseForm(); err != nil {
		var t T
		return t, err
	}

	var t T
	if err := decoder.Decode(&t, req.Form); err != nil {
		return t, err
	}
	return t, nil
}
