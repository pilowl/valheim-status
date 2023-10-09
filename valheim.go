package main

import (
	"errors"
	"time"

	"github.com/rumblefrog/go-a2s"
)

type valheimer struct {
	client *a2s.Client
}

func NewValheimer(client *a2s.Client) *valheimer {
	return &valheimer{
		client: client,
	}
}

func (v *valheimer) Close() error {
	return v.client.Close()
}

func (v *valheimer) GetPlayers() (map[ID]Player, error) {
	a2sPlayers, err := v.client.QueryPlayer()
	if err != nil {
		return nil, errors.New("failed to query a2s players")
	}

	players := make(map[ID]Player)
	for _, player := range a2sPlayers.Players {
		players[ID(player.Name)] = Player{
			Name:     player.Name,
			Playtime: time.Second * time.Duration(int(player.Duration)),
		}
	}

	return players, nil
}
