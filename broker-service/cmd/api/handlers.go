package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Err:     false,
		Message: "hit broker",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJSON(w, err)
	}
}

// to authenticate by calling authentication service then get the response back
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	data, err := json.MarshalIndent(a, "", "/t")
	if err != nil {
		app.errorJSON(w, err)
	}

	req, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(data))
	if err != nil {
		app.errorJSON(w, err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
	}
	defer resp.Body.Close()

	// make sure we get back the correct status code
	if resp.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New(("invalid credentials")))
		return
	} else if resp.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	//variable to store decoded response
	var jsonFromService jsonResponse

	err = json.NewDecoder(resp.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
	}

	if jsonFromService.Err {
		app.errorJSON(w, err, http.StatusUnauthorized)
	}

	var payload jsonResponse
	payload.Data = jsonFromService.Data
	payload.Err = false
	payload.Message = "logged in"

	app.writeJSON(w, http.StatusAccepted, payload)

}
