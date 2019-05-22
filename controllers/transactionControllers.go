package controllers

import (
	"encoding/json"
	"go-card/models"
	u "go-card/utils"
	"net/http"
)

var CreateTransaction = func(w http.ResponseWriter, r *http.Request) {
	transaction := &models.Transaction{}
	err := json.NewDecoder(r.Body).Decode(transaction)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	resp := transaction.Create()
	u.Respond(w, resp)
}

var CaptureTransaction = func(w http.ResponseWriter, r *http.Request) {
	action := &models.Action{}
	err := json.NewDecoder(r.Body).Decode(action) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := action.Capture()
	u.Respond(w, resp)
}

var ChangeTransaction = func(w http.ResponseWriter, r *http.Request) {
	action := &models.Action{}
	err := json.NewDecoder(r.Body).Decode(action)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := action.ChangeAuthorised()
	u.Respond(w, resp)
}

var RefundTransaction = func(w http.ResponseWriter, r *http.Request) {
	action := &models.Action{}
	err := json.NewDecoder(r.Body).Decode(action)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := action.Refund()
	u.Respond(w, resp)
}
