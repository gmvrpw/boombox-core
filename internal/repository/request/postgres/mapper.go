package postgres

import "gmvr.pw/boombox/pkg/model"

func requestEntityFromRequest(request *model.Request) *Request {
	return &Request{
		ID:         request.ID,
		Track:      request.Track,
		Playback:   request.Playback,
		Timecode:   request.Timecode,
		RunnerName: request.Runner.Name,
		AuthorId:   request.Author.ID,
		CreatedAt:  request.CreatedAt,
	}
}

func requestFromRequestEntity(request *Request) *model.Request {
	return &model.Request{
		ID:        request.ID,
		Track:     request.Track,
		Playback:  request.Playback,
		Timecode:  request.Timecode,
		Author:    model.User{ID: request.AuthorId},
		Runner:    model.Runner{Name: request.RunnerName, Owner: model.User{ID: request.AuthorId}},
		CreatedAt: request.CreatedAt,
	}
}
