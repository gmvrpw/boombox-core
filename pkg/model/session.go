package model

import "github.com/google/uuid"

type RunnerSessionData struct {
	Timecode uint64
	Audio    []byte
	Finished bool
}

type RunnerSession struct {
	ID uuid.UUID
	Runner
	Request
	Data <-chan *RunnerSessionData
}

type RunnerSessionNotFoundError struct{}

func (e *RunnerSessionNotFoundError) Error() string {
	return "runner session not found"
}

type RunnerSessionEmptyError struct{}

func (e *RunnerSessionEmptyError) Error() string {
	return "nothing to be piped"
}
