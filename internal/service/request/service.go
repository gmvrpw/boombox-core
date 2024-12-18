package request

import (
	"log/slog"
	"sync"

	"gmvr.pw/boombox/pkg/model"
)

var players = new(sync.Map)

type RequestService struct {
	runnerRepository  RunnerRepository
	requestRepository RequestRepository

	logger *slog.Logger
}

func NewRequestService(logger *slog.Logger) (*RequestService, error) {
	return &RequestService{logger: logger}, nil
}

func (s *RequestService) Init(
	runnerRepository RunnerRepository,
	requestRepository RequestRepository,
) {
	s.runnerRepository = runnerRepository
	s.requestRepository = requestRepository
}

func (s *RequestService) Request(req *model.Request) error {
	s.logger.Info("track requested", "user", req.Author.ID.String())
	if req.Runner.Name == "" {
		runners := s.runnerRepository.GetRunnersByTrackUrl(req.Track.Url)
		if len(runners) == 0 {
			return &model.UnplayableTrackError{}
		}
		if len(runners) > 1 {
			options := []*model.Request{}
			for index, runner := range runners {
				option := *req
				option.Runner = runner
				options[index] = &option
			}

			return &model.UnspecifiedRequestError{Options: options}
		}
		req.Runner = runners[0]
	}
	s.logger.Info("request runner specified", "runner", req.Runner)

	err := s.requestRepository.Create(req)
	if err != nil {
		return err
	}

	return nil
}

func (s *RequestService) Queue(user *model.User, pages *model.Pages) ([]*model.Request, error) {
	return s.requestRepository.GetPagedNewFirstRequestsByAuthorId(user.ID, pages)
}
