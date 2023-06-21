package models

import "errors"

var (
	NotFound      = errors.New("NotFound")
	Conflict      = errors.New("Conflict")
	InternalError = errors.New("InternalError")
)
