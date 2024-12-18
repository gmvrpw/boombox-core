package player

import (
	"log/slog"
	"sync"

	"gmvr.pw/boombox/pkg/model"
)

var players = new(sync.Map)

type PlayerService struct {
	trackRepository   TrackRepository
	requestRepository RequestRepository
	runnerRepository  RunnerRepository
	sessionRepository SessionRepository

	logger *slog.Logger
}

func NewPlayerService(logger *slog.Logger) (*PlayerService, error) {
	return &PlayerService{logger: logger}, nil
}

func (s *PlayerService) Init(
	trackRepository TrackRepository,
	requestRepository RequestRepository,
	runnerRepository RunnerRepository,
	sessionRepository SessionRepository,
) {
	s.trackRepository = trackRepository
	s.requestRepository = requestRepository
	s.runnerRepository = runnerRepository
	s.sessionRepository = sessionRepository
}

func (s *PlayerService) Play(user *model.User, target chan<- []byte) (*model.Request, error) {
	stop := make(chan bool, 1)
	ret := make(chan *model.Request, 1)
	_, exist := players.LoadOrStore(
		user.ID.String(),
		&model.Player{Stop: stop, Return: ret},
	)
	if exist {
		s.logger.Info("player already exist, skiped")
		return nil, &model.PlayerAlreadyExistsError{}
	}

	req, err := s.requestRepository.GetOldestRequestByUserId(user.ID)
	if err != nil {
		return nil, err
	}

	go func() {
		var session model.RunnerSession

		defer close(ret)
		for {
			req, err := s.requestRepository.GetOldestRequestByUserId(user.ID)
			req.Runner, err = s.runnerRepository.GetRunnerByNameAndOwnerId(
				req.Runner.Name,
				req.Runner.Owner.ID,
			)
			if err != nil {
				s.logger.Error("cannot get next track", "error", err)
				return
			}

			// TODO: need transaction
			session = model.RunnerSession{Request: *req}
			err = s.sessionRepository.Create(&session)
			if err != nil {
				s.logger.Error("cannot create runner session", "error", err)
				return
			}
			defer s.sessionRepository.Delete(&session)

			err = s.requestRepository.UpdateRequestStatusById(req.ID, model.RequestStatusRunned)
			if err != nil {
				s.logger.Error("cannot update request status", "error", err)
				return
			}

		L:
			for {
				select {
				case <-stop:
					ret <- &session.Request
					return
				case d := <-session.Data:
					if d.Finished {
						break L
					}
					req.Playback = d.Timecode
					target <- d.Audio
				}
			}

			err = s.requestRepository.UpdateRequestStatusById(req.ID, model.RequestStatusDone)
			if err != nil {
				return
			}
		}
	}()

	return req, nil
}

func (s *PlayerService) Pause(owner *model.User) (*model.Request, error) {
	stored, exist := players.Load(owner.ID.String())
	if !exist {
		s.logger.Error("player does not exist")
		return nil, &model.PlayerNotExistsError{}
	}

	player, ok := stored.(*model.Player)
	if !ok {
		s.logger.Error("cannot cast player")
		return nil, &model.PlayerNotExistsError{}
	}

	// TODO: need transaction
	player.Stop <- true
	req, ok := <-player.Return
	if !ok {
		s.logger.Error("cannot read return value")
		return nil, &model.PlayerNotExistsError{}
	}

	err := s.requestRepository.UpdateRequestStatusById(req.ID, model.RequestStatusPaused)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *PlayerService) Skip(owner *model.User) (*model.Request, error) {
	stored, exist := players.Load(owner.ID.String())
	if !exist {
		s.logger.Error("player does not exist")
		return nil, &model.PlayerNotExistsError{}
	}

	player, ok := stored.(*model.Player)
	if !ok {
		s.logger.Error("cannot cast player")
		return nil, &model.PlayerNotExistsError{}
	}

	player.Stop <- true
	req, ok := <-player.Return
	if !ok {
		s.logger.Error("cannot read return value")
		return nil, &model.PlayerNotExistsError{}
	}

	err := s.requestRepository.UpdateRequestStatusById(req.ID, model.RequestStatusSkipped)
	if err != nil {
		return nil, err
	}

	return req, nil
}
