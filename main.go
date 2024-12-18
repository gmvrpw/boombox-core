package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"

	"gmvr.pw/boombox/config"
	"gmvr.pw/boombox/internal/controller/discord"
	postgresRequestRepository "gmvr.pw/boombox/internal/repository/request/postgres"
	runtimeRunnerRepository "gmvr.pw/boombox/internal/repository/runner/runtime"
	runtimeSessionRepository "gmvr.pw/boombox/internal/repository/session/runtime"
	httpTrackRepository "gmvr.pw/boombox/internal/repository/track/http"
	"gmvr.pw/boombox/internal/service/player"
	"gmvr.pw/boombox/internal/service/request"
)

func main() {
	var err error

	logger := slog.Default()

	cfg := &config.Config{}

	path := os.Getenv("BOOMBOX_CONFIG_FILE")
	cfg, err = config.NewConfig(path)
	if err != nil {
		log.Fatalf("cannot get config. %s", err)
	}

	trackRepository, err := httpTrackRepository.NewHttpTrackRepository(logger)
	if err != nil {
		log.Fatalf("cannot create track repository. %s", err)
	}

	requestRepository, err := postgresRequestRepository.NewPostgresRequestRepository(
		cfg.Request,
		logger,
	)
	if err != nil {
		log.Fatalf("cannot create queue repository. %s", err)
	}

	runnerRepository, err := runtimeRunnerRepository.NewRuntimeRunnerRepository(cfg.Runners, logger)
	if err != nil {
		log.Fatalf("cannot create runner repository. %s", err)
	}

	sessionRepository, err := runtimeSessionRepository.NewRuntimeSessionRepository(
		cfg.Runners,
		logger,
	)
	if err != nil {
		log.Fatalf("cannot create runner repository. %s", err)
	}

	requestService, err := request.NewRequestService(logger)
	if err != nil {
		log.Fatalf("cannot create player service. %s", err)
	}

	playerService, err := player.NewPlayerService(logger)
	if err != nil {
		log.Fatalf("cannot create player service. %s", err)
	}

	discordController, err := discord.NewDiscordController(logger)
	if err != nil {
		log.Fatalf("cannot create discord controller. %s", err)
	}

	requestService.Init(runnerRepository, requestRepository)
	playerService.Init(trackRepository, requestRepository, runnerRepository, sessionRepository)
	discordController.Init(requestService, playerService)

	go func() {
		err = discordController.Serve()
		if err != nil {
			log.Fatalf("failed when serving discord. %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)
	<-stop
}
