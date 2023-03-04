package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/KishorPokharel/taskque/postgres"
	validator "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/julienschmidt/httprouter"
)

func (app *application) handleUserRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	input := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		app.errorResponse(w, http.StatusBadRequest, "Bad request body", err)
		return
	}
	if err := validator.ValidateStruct(&input,
		validator.Field(&input.Username, validator.Required, validator.Length(3, 20)),
		validator.Field(&input.Email, validator.Required, is.Email),
		validator.Field(&input.Password, validator.Required, validator.Length(10, 30)),
	); err != nil {
		out := map[string]any{
			"success": false,
			"errors":  err,
		}
		app.jsonResponse(w, http.StatusBadRequest, out)
		return
	}
	user := &postgres.User{
		Username: input.Username,
		Email:    input.Email,
	}
	err := user.Password.Set(input.Password)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	err = app.service.User.Create(user)
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrDuplicateEmail):
			out := map[string]any{
				"success": false,
				"errors": map[string]any{
					"email": "email already exists",
				},
			}
			app.jsonResponse(w, http.StatusBadRequest, out)
			return
		default:
			app.errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
			return
		}
	}
	out := map[string]any{
		"success": true,
		"message": "User registered successfully",
		"data": map[string]any{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	}
	app.jsonResponse(w, http.StatusCreated, out)
}
