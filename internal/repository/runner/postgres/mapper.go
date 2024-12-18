package postgres

import (
	"gmvr.pw/boombox/pkg/model"
)

func runnerFromRunnerEntity(entity *Runner) *model.Runner {
	return &model.Runner{
		Name:  entity.Name,
		Owner: model.User{ID: entity.OwnerId},
		Url:   entity.Url,
	}
}

func runnerEntityFromRunner(runner *model.Runner) *Runner {
	return &Runner{
		Name:    runner.Name,
		OwnerId: runner.Owner.ID,
		Url:     runner.Url,
	}
}
