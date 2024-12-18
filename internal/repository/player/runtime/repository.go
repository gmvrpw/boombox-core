package runtime

import (
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

type RuntimePlayerRepository struct {
	connections sync.Map
	logger      *slog.Logger
}

func NewRuntimeSessionRepository(
	logger *slog.Logger,
) (*RuntimePlayerRepository, error) {
	return &RuntimePlayerRepository{logger: logger}, nil
}

func GetPlyaerByOwnerID(ownerId uuid.UUID) {

}
