package postgres 

import (
	"log/slog"
	"net/http"

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

func (r *RuntimeRunnerRepository) GetRunnersByTrackUrl(url string) []model.Runner {
	return []model.Runner{}
}
