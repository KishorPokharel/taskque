package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) handleUserRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	input := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad Request Body", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, input)
}
