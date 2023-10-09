package main

import (
	"fmt"
	"time"
)

type eventer struct{}

const (
	eventLimit = 10

	eventPlayerConnected    = "Player %s has joined the server"
	eventPlayerDisconnected = "Player %s has been disconnected"
)

func NewEventer() *eventer {
	return &eventer{}
}

func (e *eventer) Analyze(pastPlayerList, presentPlayerList map[ID]Player, events []Event) []Event {
	// Get events depending on comparison of past and present events
	events = append(events, e.checkNewConnections(pastPlayerList, presentPlayerList)...)
	events = append(events, e.checkDisconnects(pastPlayerList, presentPlayerList)...)

	if len(events) >= eventLimit {
		events = events[len(events)-eventLimit:]
	}

	return events
}

func (e *eventer) checkDisconnects(pastPlayerList, presentPlayerList map[ID]Player) []Event {
	events := make([]Event, 0)
	for player := range pastPlayerList {
		if _, ok := presentPlayerList[player]; !ok {
			events = append(events, Event{
				When:        time.Now(),
				Description: fmt.Sprintf(eventPlayerDisconnected, player),
			})
		}
	}

	return events
}

func (e *eventer) checkNewConnections(pastPlayerList, presentPlayerList map[ID]Player) []Event {
	events := make([]Event, 0)
	for player := range presentPlayerList {
		if _, ok := pastPlayerList[player]; !ok {
			events = append(events, Event{
				When:        time.Now(),
				Description: fmt.Sprintf(eventPlayerConnected, player),
			})
		}
	}

	return events
}
