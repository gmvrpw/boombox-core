package discord

import "github.com/bwmarrin/discordgo"

var Commands = map[string]discordgo.ApplicationCommand{
	"play": {
		Name:        "play",
		Description: "Queue track",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "request",
				Description: "Audio request string",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
			{
				Name:        "runner",
				Description: "Name of runner to be played on",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
			{
				Name:        "channel",
				Description: "Name of the voice channel to play",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	},
	"pause": {
		Name:        "pause",
		Description: "Pause player",
	},
	"queue": {
		Name:        "queue",
		Description: "Show queue",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "page",
				Description: "Number of queue page to be displayed",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Required:    false,
			},
		},
	},
}
