package runtime

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"gmvr.pw/boombox/config"
	"gmvr.pw/boombox/pkg/model"
)

type RuntimeSessionRepository struct {
	sessions sync.Map
	client   *http.Client
	logger   *slog.Logger
}

func NewRuntimeSessionRepository(
	runners []config.RunnerConfig,
	logger *slog.Logger,
) (*RuntimeSessionRepository, error) {
	return &RuntimeSessionRepository{client: &http.Client{}, logger: logger}, nil
}

func (r *RuntimeSessionRepository) Create(session *model.RunnerSession) error {
	var err error

	entity := sessionEntityFromSession(session)
	r.logger.Info("entity created", "entity", entity, "session", session)

	entity.Port = 2001
	var con *net.UDPConn
	for {
		if entity.Port >= 3000 {
			r.logger.Error("cannot find unusing port")
			return errors.New("cannot find unusing port")
		}

		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", entity.Port))
		if err != nil {
			r.logger.Error("cannot resolve udp address", "error", err)
			return err
		}

		con, err = net.ListenUDP("udp", addr)
		if err != nil {
			entity.Port++
			continue
		}
		break
	}
	defer func() {
		if _, loaded := r.sessions.Load(session.ID.String()); !loaded {
			con.Close()
		}
	}()

	d := make(chan *model.RunnerSessionData)
	stop := make(chan bool)
	go func() {
		m := make([]byte, 1024)

		con.SetDeadline(time.Now().Add(time.Second * 20))
		for {
			select {
			case <-stop:
				return
			default:
				n, _, err := con.ReadFromUDP(m)
				if err != nil {
					r.logger.Error("cannot read from udp", "error", err)
					return
				}
				if n < 4 {
					r.logger.Error("unexpected datagram size", "size", n)
					return
				}

				timecode := binary.BigEndian.Uint64(m[:8])
				d <- &model.RunnerSessionData{Timecode: timecode, Audio: m[8:], Finished: timecode == ^uint64(0)}

				con.SetDeadline(time.Now().Add(time.Second * 1))
			}
		}
	}()

	session.Data = d
	entity.Stop = stop

	data, err := json.Marshal(entity)
	if err != nil {
		r.logger.Error("cannot marshal body", "error", err)
		return err
	}
	r.logger.Info("body marshalled", "body", entity)

	path, err := url.JoinPath(session.Request.Runner.Url, "./runners")
	if err != nil {
		r.logger.Error("cannot create request path", "error", err)
		return err
	}

	res, err := r.client.Post(path, "application/json", bytes.NewBuffer(data))
	if err != nil {
		r.logger.Error("cannot send request", "path", path, "error", err)
		return err
	}

	if res.StatusCode != http.StatusOK {
		r.logger.Error("bad request", "status_code", res.StatusCode)
		return errors.New("bad request")
	}

	err = json.NewDecoder(res.Body).Decode(&entity)
	if err != nil {
		return err
	}

	session.ID, err = uuid.Parse(entity.ID)
	if err != nil {
		return err
	}

	r.sessions.Store(entity.ID, &entity)
	return nil
}

func (r *RuntimeSessionRepository) Delete(session *model.RunnerSession) error {
	stored, loaded := r.sessions.LoadAndDelete(session.ID.String())
	if !loaded {
		r.logger.Error("cannot find session", "id", session.ID.String())
		return &model.RunnerSessionNotFoundError{}
	}

	entity, ok := stored.(*RunnerSession)
	if !ok {
		r.logger.Error("cannot find session", "id", session.ID.String())
		return &model.RunnerSessionNotFoundError{}
	}

	path, err := url.JoinPath(session.Request.Runner.Url, "runners", session.ID.String())
	if err != nil {
		r.logger.Error("cannot join path", "error", err)
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, path, bytes.NewBuffer([]byte{}))
	if err != nil {
		r.logger.Error("cannot create request", "error", err)
		return err
	}

	res, err := r.client.Do(req)
	if err != nil {
		r.logger.Error("cannot send request", "error", err)
		return err
	}

	if res.StatusCode != http.StatusOK {
		r.logger.Error("bad request", "code", res.StatusCode)
		return errors.New("bad request")
	}

	entity.Stop <- true
	return nil
}
