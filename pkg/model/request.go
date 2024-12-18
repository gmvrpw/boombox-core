package model

import (
	"time"

	"github.com/google/uuid"
)

type RequestStatus = string

var (
	RequestStatusQueued  = "queued"
	RequestStatusRunned  = "runned"
	RequestStatusPaused  = "paused"
	RequestStatusSkipped = "skipped"
	RequestStatusDone    = "done"
)

type Request struct {
	ID uuid.UUID
	Track
	Runner    Runner
	Playback  uint64
	Timecode  uint64
	Author    User
	Status    RequestStatus
	CreatedAt time.Time
}

type RequestNotFoundError struct{}

func (e *RequestNotFoundError) Error() string {
	return "request not found"
}
