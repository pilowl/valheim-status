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
	EnvStatusUpdateFrequency = "VALHEIM_DISCORD_STATUS_UPDATE_RATE"
	EnvServerIP              = "VALHEIM_DISCORD_STATUS_STEAM_STATUS_IP"
)

type config struct {
	DiscordAPIKey         string
	StatusUpdateFrequency time.Duration
	ServerIP              string
}

func (c *config) String() string {
	s := ""
	s += fmt.Sprintf("%s=%v;", EnvDiscordAPIKey, c.DiscordAPIKey)
	s += fmt.Sprintf("%s=%v;", EnvStatusUpdateFrequency, c.StatusUpdateFrequency)
	s += fmt.Sprintf("%s=%v.", EnvServerIP, c.ServerIP)
	return s
}

func NewMockConfig() (*config, error) {
	return &config{
		DiscordAPIKey:         "123",
		StatusUpdateFrequency: 60 * time.Second,
		ServerIP:              "11.22.33.44",
	}, nil
}

func NewConfig() (*config, error) {
	emptyEnvVariableError := func(variable string) error {
		return fmt.Errorf("%s environment variable is not set", variable)
	}

	discordAPIKey := os.Getenv(EnvDiscordAPIKey)
	if len(discordAPIKey) == 0 {
		return nil, emptyEnvVariableError(EnvDiscordAPIKey)
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

	return &config{
		DiscordAPIKey:         discordAPIKey,
		StatusUpdateFrequency: time.Second * time.Duration(statusUpdateFrequency),
		ServerIP:              serverIP,
	}, nil
}
