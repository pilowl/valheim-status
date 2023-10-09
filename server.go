package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/rumblefrog/go-a2s"
	"golang.org/x/exp/slog"
)

type DiscordSession interface {
	Open() error
	SendStatusMessage(channelID string, content string) error
	Close() error
}

type Valheimer interface {
	GetPlayers() (map[ID]Player, error)
	Close() error
}

type Eventer interface {
	Analyze(pastPlayerList, presentPlayerList map[ID]Player, events []Event) []Event
}

type Messenger interface {
	BuildMessage(serverStatus Status) (string, error)
}

func run(logger *slog.Logger, config *Config) error {
	var (
		discord   DiscordSession
		valheim   Valheimer
		eventer   Eventer
		messenger Messenger
	)

	discordSession, err := discordgo.New("Bot " + config.DiscordAPIKey)
	if err != nil {
		return errors.Wrap(err, "failed to initialize discord client")
	}
	discord = NewDiscorder(discordSession, config.DiscordBotID)
	if err := discord.Open(); err != nil {
		return errors.Wrap(err, "open discord session")
	}

	defer func() {
		if err := discord.Close(); err != nil {
			logger.With(err).Info("Failed to close discord connections")
		}
	}()

	a2sClient, err := a2s.NewClient(fmt.Sprintf("%s:%s", config.ServerIP, "2457"),
		a2s.SetMaxPacketSize(14000),
		a2s.TimeoutOption(time.Second*5),
	)
	if err != nil {
		logger.With("error", err).Error("failed to REinitialize a2s")
		return nil
	}
	valheim = NewValheimer(a2sClient)
	defer func() {
		if err := valheim.Close(); err != nil {
			logger.With(err).Error("Failed to close valheim status server client")
		}
	}()

	tabler := NewTabler()
	messenger = NewMessenger(tabler)

	eventer = NewEventer()

	updateRate := config.StatusUpdateFrequency

	status := Status{
		Players: make(map[ID]Player),
		Events:  make([]Event, 0),
	}

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt)

	t := time.Tick(updateRate)
	for {
		select {
		case <-t:
			players, err := valheim.GetPlayers()
			if err != nil {
				logger.With("error", err).Error("failed to get valheim server players, reinitializing a2s client now...")

				a2sClient = reinitializeA2S(logger, config.ServerIP)
				if a2sClient != nil {
					valheim = NewValheimer(a2sClient)
				}
				continue
			}

			status.Events = eventer.Analyze(status.Players, players, status.Events)
			status.Players = players

			message, err := messenger.BuildMessage(status)
			if err != nil {
				logger.With("error", err).Error("failed to build a status message")
				continue
			}

			if err := discord.SendStatusMessage(config.ChannelID, message); err != nil {
				logger.With("Error", err).Error("failed to send message to the discord server")
			}
		case <-stopCh:
			return nil
		}
	}
}

func reinitializeA2S(logger *slog.Logger, ip string) *a2s.Client {
	a2sClient, err := a2s.NewClient(fmt.Sprintf("%s:%s", ip, "2457"),
		a2s.SetMaxPacketSize(14000),
		a2s.TimeoutOption(time.Second*5),
	)
	if err != nil {
		logger.With("error", err).Error("failed to REinitialize a2s")
		return nil
	}

	return a2sClient
}
