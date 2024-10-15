package main

import (
	"authentication/data"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	//validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.errorJson(w, errors.New("Invalid credentials 1"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)

	if err != nil || !valid {
		app.errorJson(w, errors.New("Invalid credentials 2"), http.StatusBadRequest)
	}

	payload := jsonReponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) Register(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Active    bool   `json:"active"`
	}

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	var newUser data.User
	newUser.Email = requestPayload.Email
	newUser.FirstName = requestPayload.FirstName
	newUser.LastName = requestPayload.LastName
	newUser.Password = requestPayload.Password
	newUser.Active = requestPayload.Active

	userID, err := app.Models.User.Insert(newUser)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err) // Log chi tiết lỗi
		app.errorJson(w, errors.New("Unable to insert user into database"), http.StatusInternalServerError)
		return
	}

	user, err := app.Models.User.GetOne(userID)
	if err != nil {
		app.errorJson(w, errors.New("User not found after insert"), http.StatusInternalServerError)
		return
	}

	payload := jsonReponse{
		Error:   false,
		Message: fmt.Sprintf("User %s successfully registered", user.Email),
		Data:    user,
	}

	app.writeJson(w, http.StatusAccepted, payload)
}
