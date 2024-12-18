package discord

import (
	"github.com/bwmarrin/discordgo"
	"gmvr.pw/boombox/pkg/model"
)

func (c *DiscordController) Pause(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) (*discordgo.InteractionResponseData, error) {
	var err error

	req, err := c.playerService.Pause(&model.User{ID: id})
	if err != nil {
		if _, ok := err.(*model.PlayerNotExistsError); ok {
			// TODO: response about it
		}

		return nil, err
	}

	return queuedFromRequest(req), nil
}
