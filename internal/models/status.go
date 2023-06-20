package models

// easyjson -all ./internal/models/status.go

type Status struct {
	Post   int `json:"post"`
	Author int `json:"author"`
	Forum  int `json:"forum"`
	Thread int `json:"thread"`
}
