package main

import "net/http"

type WebhookPayload struct {
	Event string `json:"event"`
	Data  User   `json:"data"`
}

func (cfg *apiConfig) handleWebhook(w http.ResponseWriter, r *http.Request) {
	//
}
