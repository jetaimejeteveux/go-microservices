package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Err     bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

/*
To summarize, this Go function is designed to handle the reading and decoding of JSON data from an HTTP request.
It ensures that the request body is limited in size, decodes the JSON data, and checks that there's only one JSON value present.
If everything succeeds, it returns nil to indicate success. Otherwise, it returns an error.
*/
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxbytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxbytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contains 1 single json value")
	}
	return nil

}

/*
this Go function is designed to serialize data into JSON format, set appropriate headers, and send an HTTP response.
It handles cases where optional headers are provided, and it returns an error if any issues occur during the process.
*/
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil

}

/*
To summarize, this Go function errorJSON is a utility for sending error responses in JSON format.
It sets the HTTP status code based on the provided status or defaults to http.StatusBadRequest
It creates a JSON payload indicating an error, including the error message obtained from the err parameter.
Finally, it calls another function (writeJSON) to send the JSON response.
*/
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	var payload jsonResponse
	payload.Err = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)

}
