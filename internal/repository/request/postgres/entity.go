package postgres

import (
	"time"

	"github.com/google/uuid"
	"gmvr.pw/boombox/pkg/model"
)

type Request struct {
	ID         uuid.UUID   `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Track      model.Track `gorm:"embedded;embeddedPrefix:track_"`
	Playback   uint64
	Timecode   uint64
	RunnerName string
	AuthorId   uuid.UUID `gorm:"type:uuid"`
	Status     model.RequestStatus
	CreatedAt  time.Time
}
