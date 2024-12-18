package discord

import (
	"github.com/bwmarrin/discordgo"
	"gmvr.pw/boombox/pkg/model"
)

func (c *DiscordController) Queue(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) (*discordgo.InteractionResponseData, error) {
	page := pageFromCommandData(i.ApplicationCommandData())
	pages := model.Pages{Size: 10, Start: page, Stop: page + 1}

	req, _ := c.requestService.Queue(
		&model.User{ID: id},
		&pages,
	)

	return queueFromRequests(req, pages.Start*pages.Size+1), nil
}
