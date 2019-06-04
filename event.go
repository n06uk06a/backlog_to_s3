package main

import (
	"encoding/json"
	"log"
	"net/url"
)

// BacklogGitWebhookEvent struct of Backlog git webhook event
type BacklogGitWebhookEvent struct {
	Before     string `json:"before"`
	After      string `json:"after"`
	Ref        string `json:"ref"`
	Repository struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"repository"`
}

// NewBacklogGitWebhookEvent initialize BacklogGitWebhookEvent from request body
func NewBacklogGitWebhookEvent(body string) *BacklogGitWebhookEvent {
	values, err := url.ParseQuery(body)
	if err != nil {
		log.Fatal(err)
	}
	payload := values["payload"][0]
	var event = new(BacklogGitWebhookEvent)
	if err := json.Unmarshal([]byte(payload), event); err != nil {
		log.Fatal(err)
	}
	return event
}
