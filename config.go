package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	EnvDiscordAPIKey         = "VALHEIM_DISCORD_STATUS_API_KEY"
	EnvDiscordBotID          = "VALHEIM_DISCORD_BOT_ID"
	EnvDiscordChannelID      = "VALHEIM_DISCORD_CHANNEL_ID"
	EnvStatusUpdateFrequency = "VALHEIM_DISCORD_STATUS_UPDATE_RATE"
	EnvServerIP              = "VALHEIM_DISCORD_STATUS_STEAM_STATUS_IP"
)

type Config struct {
	DiscordBotID          string
	DiscordAPIKey         string
	ChannelID             string
	StatusUpdateFrequency time.Duration
	ServerIP              string
}

func (c *Config) String() string {
	s := ""
	s += fmt.Sprintf("%s=%v;", EnvDiscordAPIKey, c.DiscordAPIKey)
	s += fmt.Sprintf("%s=%v;", EnvDiscordBotID, c.DiscordBotID)
	s += fmt.Sprintf("%s=%v;", EnvDiscordChannelID, c.ChannelID)
	s += fmt.Sprintf("%s=%v;", EnvStatusUpdateFrequency, c.StatusUpdateFrequency)
	s += fmt.Sprintf("%s=%v.", EnvServerIP, c.ServerIP)
	return s
}

func NewConfig() (*Config, error) {
	emptyEnvVariableError := func(variable string) error {
		return fmt.Errorf("%s environment variable is not set", variable)
	}

	discordAPIKey := os.Getenv(EnvDiscordAPIKey)
	if len(discordAPIKey) == 0 {
		return nil, emptyEnvVariableError(EnvDiscordAPIKey)
	}

	channelID := os.Getenv(EnvDiscordChannelID)
	if len(channelID) == 0 {
		return nil, emptyEnvVariableError(EnvDiscordChannelID)
	}

	statusUpdateFrequency, err := strconv.Atoi(os.Getenv(EnvStatusUpdateFrequency))
	if err != nil {
		if len(os.Getenv(EnvStatusUpdateFrequency)) == 0 {
			return nil, emptyEnvVariableError(EnvStatusUpdateFrequency)
		}

		return nil, errors.New("unable to parse update frequency. Needs to be integer, represeting seconds.")
	}

	serverIP := os.Getenv(EnvServerIP)
	if len(serverIP) == 0 {
		return nil, emptyEnvVariableError(EnvServerIP)
	}

	return &Config{
		DiscordAPIKey:         discordAPIKey,
		StatusUpdateFrequency: time.Second * time.Duration(statusUpdateFrequency),
		ServerIP:              serverIP,
	}, nil
}
