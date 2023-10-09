package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type discorder struct {
	session *discordgo.Session
	botID   string
}

func NewDiscorder(session *discordgo.Session, botID string) *discorder {
	return &discorder{
		botID:   botID,
		session: session,
	}
}

func (d *discorder) Open() error {
	if err := d.session.Open(); err != nil {
		return errors.Wrap(err, "open discord session")
	}

	return nil
}

func (d *discorder) SendStatusMessage(channelID string, content string) error {
	messages, err := d.session.ChannelMessages(channelID, 100, "", "", "")
	if err != nil {
		return errors.Wrap(err, "read channel messages")
	}

	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if msg.Author.ID == d.botID {
			if _, err := d.session.ChannelMessageEdit(channelID, msg.ID, content); err != nil {
				return errors.Wrap(err, "edit bot channel message")
			}
		}

		return nil
	}

	if _, err := d.session.ChannelMessageSend(channelID, content); err != nil {
		return errors.Wrap(err, "send message to discord channel")
	}

	return nil
}

func (d *discorder) Close() error {
	if err := d.session.Close(); err != nil {
		return errors.Wrap(err, "close discord session")
	}

	return nil
}
