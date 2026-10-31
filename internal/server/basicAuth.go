package server

import (
	"errors"
	"net/http"
	"webtest/internal/logger"
	"webtest/internal/models"
)

var (
	//ERRORS
	AuthErr = errors.New("Аутентификация не пройдена!")
)

func requierAuth(w http.ResponseWriter, r *http.Request) error {

	user, pass, ok := r.BasicAuth()

	if !ok || user != username || pass != password {
		w.Header().Set("WWW-Authenticate", `Basic realm="Cookbook Area"`)

		go logger.LogPrint(AuthErr.Error(), models.Error)

		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return AuthErr
	}

	return nil
}
