package postgres

import "github.com/google/uuid"

type Runner struct {
	Name    string    ``
	OwnerId uuid.UUID ``
	Url     string    `json:"url"`
}
