package discord

import (
	"gmvr.pw/boombox/pkg/model"
)

type RequestService interface {
	Request(req *model.Request) error
	Queue(author *model.User, pages *model.Pages) ([]*model.Request, error)
}

type PlayerService interface {
	Play(user *model.User, target chan<- []byte) (*model.Request, error)
	Pause(user *model.User) (*model.Request, error)
}
