package models

// easyjson -all ./internal/models/status.go

type Status struct {
	User   int `json:"user"`
	Forum  int `json:"forum"`
	Thread int `json:"thread"`
	Post   int `json:"post"`
}
