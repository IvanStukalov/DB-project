package models

import "time"

// easyjson -all ./internal/models/post.go

type Post struct {
	ID       int       `json:"id"`
	Parent   int       `json:"parent,omitempty"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Forum    string    `json:"forum"`
	Thread   int       `json:"thread"`
	Created  time.Time `json:"created,omitempty"`
}
