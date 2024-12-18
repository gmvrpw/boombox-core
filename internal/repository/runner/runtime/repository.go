package runtime

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"gmvr.pw/boombox/config"
	"gmvr.pw/boombox/pkg/model"
)

type RuntimeRunnerRepository struct {
	runners []config.RunnerConfig
	client  *http.Client
	logger  *slog.Logger
}

func NewRuntimeRunnerRepository(
	runners []config.RunnerConfig,
	logger *slog.Logger,
) (*RuntimeRunnerRepository, error) {
	return &RuntimeRunnerRepository{runners: runners, client: &http.Client{}, logger: logger}, nil
}

func (r *RuntimeRunnerRepository) GetRunnerByNameAndOwnerId(
	name string,
	ownerId uuid.UUID,
) (model.Runner, error) {
	for _, entity := range r.runners {
		if entity.Name == name && (entity.Owner == uuid.Nil || entity.Owner == ownerId) {
			return *runnerFromRunnerEntity(&entity), nil
		}
	}

	return model.Runner{}, &model.RunnerSessionNotFoundError{}
}

func (r *RuntimeRunnerRepository) GetRunnersByTrackUrl(url string) []model.Runner {
	runners := []model.Runner{}

	for _, entity := range r.runners {
		runners = append(runners, *runnerFromRunnerEntity(&entity))
	}

	return runners
}
