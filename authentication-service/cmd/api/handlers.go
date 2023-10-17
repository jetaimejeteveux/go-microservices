package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
	}
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("bad credential"), http.StatusBadRequest)
	}
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("bad credential"), http.StatusBadRequest)
	}
	payload := jsonResponse{
		Err:     false,
		Message: fmt.Sprintf("logged in, user = %s", user.Email),
		Data:    user,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}
