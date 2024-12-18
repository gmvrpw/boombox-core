package discord

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"gmvr.pw/boombox/pkg/model"
)

var tracksStore = sync.Map{}

var id = uuid.New()

func (c *DiscordController) Play(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) (*discordgo.InteractionResponseData, error) {
	var err error

	target, err := targetFromRequest(s, i)
	if err != nil {
		return nil, err
	}

	req := requestFromCommandData(i.ApplicationCommandData())
	req.Author.ID = id

	err = c.requestService.Request(req)
	if err != nil {
		if _, ok := err.(*model.UnplayableTrackError); ok {
			// TODO: response about it
		}

		if e, ok := err.(*model.UnspecifiedRequestError); ok {
			return specifyFromRequests(e.Options, i.ID), nil
		}

		return nil, err
	}

	_, err = c.playerService.Play(&req.Author, target)
	if err != nil {
		return nil, err
	}

	return queuedFromRequest(req), nil
}

func (c *DiscordController) TrackSpecified(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) (*discordgo.InteractionResponseData, error) {
	var err error

	target, err := targetFromRequest(s, i)
	if err != nil {
		return nil, err
	}

	data := i.MessageComponentData()
	req, err := requestFromComponentData(&data)
	if err != nil {
		if e, ok := err.(*model.UnspecifiedRequestError); ok {
			return specifyFromRequests(e.Options, i.ID), nil
		}
		return nil, err
	}

	err = c.requestService.Request(req)
	if err != nil {
		if _, ok := err.(*model.UnplayableTrackError); ok {
			// TODO: response about it
		}

		if e, ok := err.(*model.UnspecifiedRequestError); ok {
			return specifyFromRequests(e.Options, i.ID), nil
		}

		return nil, err
	}

	_, err = c.playerService.Play(&req.Author, target)
	if err != nil {
		return nil, err
	}

	return queuedFromRequest(req), nil
}
