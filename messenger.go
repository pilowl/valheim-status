package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"
)

type Tabler interface {
	Build(headers []string, rows [][]string) (string, error)
}

type messenger struct {
	tabler Tabler
}

func NewMessenger(tabler Tabler) *messenger {
	return &messenger{
		tabler: tabler,
	}
}

func (m *messenger) BuildMessage(serverStatus Status) (string, error) {
	playerTable, err := m.buildPlayersTable(serverStatus.Players)
	if err != nil {
		return "", errors.Wrap(err, "build players table")
	}
	eventTable, err := m.buildEventsTable(serverStatus.Events)
	if err != nil {
		return "", errors.Wrap(err, "build events table")
	}
	text := fmt.Sprintf(`
# Valheim Server
## Player count: *%d/10*
### Current players:
`, len(serverStatus.Players))
	text += "```\n"
	text += playerTable
	text += "```\n"
	text += "### Connection log: \n"
	text += "```\n"
	text += eventTable
	text += "```\n"
	text += fmt.Sprintf("### Last update: %s", time.Now().Format("2006-01-02 15:04"))

	return text, nil
}

func (m *messenger) buildEventsTable(events []Event) (string, error) {
	headers := []string{"Time", "Event"}
	columns := make([][]string, len(events))

	for idx, event := range events {
		columns[idx] = []string{event.When.Format("2006-01-02 15:04"), event.Description}
	}

	return m.tabler.Build(headers, columns)
}

func (m *messenger) buildPlayersTable(players map[ID]Player) (string, error) {
	headers := []string{"Player Name", "Playing time"}
	columns := make([][]string, len(players))

	playerSlice := make([]Player, 0, len(players))
	for _, player := range players {
		playerSlice = append(playerSlice, player)
	}

	sort.Slice(playerSlice, func(i, j int) bool {
		return playerSlice[i].Playtime > (playerSlice[j].Playtime)
	})

	for idx, player := range playerSlice {
		formattedTime := ""
		if player.Playtime.Hours() != 0 {
			formattedTime += fmt.Sprintf("%dh", int(player.Playtime.Hours()))
		}
		formattedTime += fmt.Sprintf("%dm", int(player.Playtime.Minutes()))
		columns[idx] = []string{player.Name, fmt.Sprintf("%s", formattedTime)}
	}

	return m.tabler.Build(headers, columns)
}
