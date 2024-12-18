package postgres

import (
	"errors"
	"log/slog"
	"net/url"

	"github.com/google/uuid"
	"gmvr.pw/boombox/config"
	"gmvr.pw/boombox/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRequestRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewDSN(u string) (*url.URL, error) {
	dsn, err := url.Parse(u)
	if err != nil {
		return nil, errors.New("cannot parse postgres url")
	}

	if dsn.Scheme != "" && dsn.Scheme != "postgresql" {
		return nil, errors.New("cannot use not posgresql scheme")
	}
	dsn.Scheme = "postgresql"

	user := dsn.User.Username()
	if user == "" {
		user = config.GetSecret("BOOMBOX_REQUESTS_POSTGRES_USER")
	}
	password, _ := dsn.User.Password()
	if password == "" {
		password = config.GetSecret("BOOMBOX_REQUESTS_POSTGRES_PASSWORD")
	}
	dsn.User = url.UserPassword(user, password)

	if dsn.Path == "" {
		dsn.Path = config.GetSecret("BOOMBOX_REQUESTS_POSTGRES_DB")
	}

	return dsn, nil
}

func NewPostgresRequestRepository(
	config config.RequestsConfig,
	logger *slog.Logger,
) (*PostgresRequestRepository, error) {
	var err error

	r := PostgresRequestRepository{logger: logger}

	dsn, err := NewDSN(config.Url)
	if err != nil {
		return nil, err
	}

	r.db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn.String(),
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	r.db.AutoMigrate(&Request{})

	return &r, nil
}

func (r *PostgresRequestRepository) Create(req *model.Request) error {
	entity := requestEntityFromRequest(req)

	err := r.db.Create(entity).Error
	if err != nil {
		return nil
	}

	req.ID = entity.ID
	return nil
}

func (r *PostgresRequestRepository) GetOldestRequestByUserId(
	userId uuid.UUID,
) (*model.Request, error) {
	entity := Request{}

	err := r.db.Where(&Request{AuthorId: userId}).Order("created_at asc").First(&entity).Error
	if err != nil {
		return nil, err
	}
	if entity.ID == uuid.Nil {
		return nil, &model.RequestNotFoundError{}
	}

	return requestFromRequestEntity(&entity), nil
}

func (r *PostgresRequestRepository) GetPagedNewFirstRequestsByAuthorId(
	authorId uuid.UUID,
	pages *model.Pages,
) ([]*model.Request, error) {
	entities := []Request{}

	err := r.db.Where(&Request{AuthorId: authorId}).
		Order("created_at desc").
		Offset(pages.Size*pages.Start - 1).
		Limit(pages.Size * (pages.Stop - pages.Start)).Find(&entities).Error

	requests := []*model.Request{}
	for _, entity := range entities {
		requests = append(requests, requestFromRequestEntity(&entity))
	}

	return requests, err
}

func (r *PostgresRequestRepository) UpdateRequestStatusById(
	id uuid.UUID,
	status model.RequestStatus,
) error {
	err := r.db.Model(&Request{ID: id}).Updates(&Request{Status: status}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRequestRepository) UpdateRequestPlaybackById(id uuid.UUID, playback uint64) error {
	err := r.db.Model(&Request{ID: id}).Updates(&Request{Playback: playback}).Error
	if err != nil {
		return err
	}

	return nil
}
