package discord

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"gmvr.pw/boombox/config"
)

type DiscordController struct {
	session *discordgo.Session

	requestService RequestService
	playerService  PlayerService

	logger *slog.Logger
}

func NewDiscordController(
	logger *slog.Logger,
) (*DiscordController, error) {
	var err error

	c := DiscordController{logger: logger.With("controller", "discord")}

	token := config.GetSecret("BOOMBOX_DISCORD_BOT_TOKEN")
	c.session, err = discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *DiscordController) Init(requestService RequestService, playerService PlayerService) {
	c.requestService = requestService
	c.playerService = playerService
}

func (c *DiscordController) Serve() error {
	var err error

	err = c.session.Open()
	if err != nil {
		return err
	}
	defer c.session.Close()

	// Register commands
	for name, cmd := range Commands {
		_, err = c.session.ApplicationCommandCreate(c.session.State.User.ID, "", &cmd)
		if err != nil {
			c.logger.Error("cannot create command", "name", name, "error", err)
		}
		c.logger.Info("command created", "name", name)
	}

	// Register commands handler
	c.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var (
			err  error
			data *discordgo.InteractionResponseData
		)

		c.logger.Info("interation handled")

		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			switch i.ApplicationCommandData().Name {
			case "play":
				data, err = c.Play(s, i)
			case "pause":
				data, err = c.Pause(s, i)
			case "queue":
				data, err = c.Queue(s, i)
			}

			if err != nil {
				c.logger.Error("failed to execute", "error", err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "An internal error occurred. Please check the logs for more information.",
					},
				})
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: data,
			})

		case discordgo.InteractionMessageComponent:
			c.TrackSpecified(s, i)
		}
	})

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)
	<-stop

	return nil
}
