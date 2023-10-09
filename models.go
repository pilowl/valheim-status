package main

import "time"

type ID string

type Player struct {
	Name     string
	Playtime time.Duration
}

type Event struct {
	When        time.Time
	Description string
}

type Status struct {
	Players map[ID]Player

	// First In Last Out
	Events []Event
}
