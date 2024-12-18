package runtime

import (
	"gmvr.pw/boombox/pkg/model"
)

func sessionEntityFromSession(session *model.RunnerSession) *RunnerSession {
	return &RunnerSession{
		ID:       session.ID.String(),
		Url:      session.Request.Url,
		Playback: session.Request.Playback,
		Port:     2000,
	}
}
