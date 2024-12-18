package http

import (
	"log/slog"

	"gmvr.pw/boombox/pkg/model"
)

type HttpTrackRepository struct {
	logger *slog.Logger
}

func NewHttpTrackRepository(logger *slog.Logger) (*HttpTrackRepository, error) {
	return &HttpTrackRepository{logger: logger}, nil
}

func (r *HttpTrackRepository) Search(query string) []model.Track {
	return []model.Track{}
}
