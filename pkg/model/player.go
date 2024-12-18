package model

type Player struct {
	Return <-chan *Request
	Stop   chan<- bool
}

type PlayerAlreadyExistsError struct{}

func (e *PlayerAlreadyExistsError) Error() string {
	return "player already exist"
}

type PlayerNotExistsError struct{}

func (e *PlayerNotExistsError) Error() string {
	return "player does not exist"
}
