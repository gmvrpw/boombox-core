package runtime

import (
	"gmvr.pw/boombox/config"
	"gmvr.pw/boombox/pkg/model"
)

func runnerFromRunnerEntity(entity *config.RunnerConfig) *model.Runner {
	return &model.Runner{
		Name:  entity.Name,
		Owner: model.User{ID: entity.Owner},
		Url:   entity.Url,
		Test:  entity.Test,
	}
}
