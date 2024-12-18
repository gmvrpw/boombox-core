package discord

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gmvr.pw/boombox/pkg/model"
)

func targetFromRequest(s *discordgo.Session, i *discordgo.InteractionCreate) (chan []byte, error) {
	con, connected := s.VoiceConnections[i.GuildID]
	if !connected {
		vs, err := s.State.VoiceState(i.GuildID, i.Member.User.ID)
		if err != nil {
			return nil, err
		}
		con, err = s.ChannelVoiceJoin(i.GuildID, vs.ChannelID, false, false)
		if err != nil {
			return nil, err
		}
	}

	return con.OpusSend, nil
}

func pageFromCommandData(data discordgo.ApplicationCommandInteractionData) int {
	for _, option := range data.Options {
		switch option.Name {
		case "page":
			return int(option.IntValue())
		}
	}

	return 0
}

func requestFromCommandData(data discordgo.ApplicationCommandInteractionData) *model.Request {
	req := &model.Request{
		Author: model.User{ID: id},
	}

	for _, option := range data.Options {
		switch option.Name {
		case "request":
			req.Url = option.StringValue()
		case "runner":
			req.Runner.Name = option.StringValue()
		}
	}

	return req
}

func requestFromComponentData(
	data *discordgo.MessageComponentInteractionData,
) (*model.Request, error) {
	stored, ok := tracksStore.LoadAndDelete(data.CustomID)
	if !ok {
		return nil, errors.New("cannot find tracks in store")
	}
	tracks := stored.([]model.Request)

	chosen, err := strconv.Atoi(data.Values[0])
	if err != nil {
		return nil, err
	}

	if chosen < 0 || chosen > len(tracks) {
		return nil, errors.New("wrong track index")
	}

	return &tracks[chosen], nil
}

func titleFromRequest(req *model.Request) string {
	label := req.Url
	if req.Track.Name != "" {
		label = req.Track.Name
	}
	if req.Track.Name != "" && req.Track.Author.Name != "" {
		label = fmt.Sprintf(
			"%s - %s",
			req.Track.Name,
			req.Track.Author,
		)
	}
	return label
}

func specifyFromRequests(reqs []*model.Request, id string) *discordgo.InteractionResponseData {
	options := []discordgo.SelectMenuOption{}
	for index, req := range reqs {
		options = append(
			options,
			discordgo.SelectMenuOption{
				Emoji:       &discordgo.ComponentEmoji{ID: "1308840582705319956"},
				Label:       titleFromRequest(req),
				Description: fmt.Sprintf("%s, %s", req.Track.Service.Name, req.Runner.Name),
				Value:       strconv.Itoa(index),
			},
		)
	}

	tracksStore.Store(id, reqs)
	return &discordgo.InteractionResponseData{
		CustomID: "track",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    id,
					MenuType:    discordgo.StringSelectMenu,
					Placeholder: "Select track to be played",
					MaxValues:   1,
					Options:     options,
				},
			}},
		}}
}

func queueFromRequests(reqs []*model.Request, base int) *discordgo.InteractionResponseData {
	var content []string
	for index, req := range reqs {
		content = append(content, fmt.Sprintf("%d. %s", base+index, titleFromRequest(req)))
	}

	return &discordgo.InteractionResponseData{
		Content: strings.Join(content, "\n"),
	}
}

func queuedFromRequest(req *model.Request) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Author: &discordgo.MessageEmbedAuthor{
					Name:    req.Track.Service.Name,
					URL:     req.Track.Service.Url,
					IconURL: req.Track.Service.Icon,
				},
				Title: titleFromRequest(req),
				URL:   req.Track.Url,
				Footer: &discordgo.MessageEmbedFooter{
					Text: req.Runner.Name,
				},
			},
		},
	}
}
