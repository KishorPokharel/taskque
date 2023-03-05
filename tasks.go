package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KishorPokharel/taskque/postgres"
	validator "github.com/go-ozzo/ozzo-validation/v4"
)

func (app *application) handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	input := struct {
		Content string `json:"content"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		app.errorResponse(
			w,
			http.StatusBadRequest,
			"Bad request body",
			fmt.Errorf("error: decoding json: %w", err),
		)
		return
	}
	if err := validator.ValidateStruct(&input,
		validator.Field(&input.Content, validator.Required),
	); err != nil {
		out := map[string]any{
			"success": false,
			"errors":  err,
		}
		app.jsonResponse(w, http.StatusBadRequest, out)
		return
	}
	user := app.contextGetUser(r)
	task := &postgres.Task{
		UserID:  user.ID,
		Content: input.Content,
	}
	if err := app.service.Task.Insert(task); err != nil {
		app.errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	out := map[string]any{
		"success": true,
		"message": "Task added successfully",
		"data": map[string]any{
			"id":         task.ID,
			"content":    task.Content,
			"created_at": task.CreatedAt,
		},
	}
	app.jsonResponse(w, http.StatusCreated, out)
}
