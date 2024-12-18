package player

import (
	"github.com/google/uuid"
	"gmvr.pw/boombox/pkg/model"
)

type TrackRepository interface {
	Search(query string) []model.Track
}

type RequestRepository interface {
	GetOldestRequestByUserId(userId uuid.UUID) (*model.Request, error)

	UpdateRequestStatusById(id uuid.UUID, status model.RequestStatus) error
	UpdateRequestPlaybackById(id uuid.UUID, timestamp uint64) error
}

type SessionRepository interface {
	Create(session *model.RunnerSession) error
	Delete(session *model.RunnerSession) error
}

type RunnerRepository interface {
	GetRunnerByNameAndOwnerId(name string, ownerId uuid.UUID) (model.Runner, error)
	GetRunnersByTrackUrl(url string) []model.Runner
}
