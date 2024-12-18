package request

import (
	"github.com/google/uuid"
	"gmvr.pw/boombox/pkg/model"
)

type RunnerRepository interface {
	GetRunnerByNameAndOwnerId(name string, ownerId uuid.UUID) (model.Runner, error)
	GetRunnersByTrackUrl(url string) []model.Runner
}

type RequestRepository interface {
	Create(req *model.Request) error

	GetPagedNewFirstRequestsByAuthorId(
		authorId uuid.UUID,
		pages *model.Pages,
	) ([]*model.Request, error)
}
