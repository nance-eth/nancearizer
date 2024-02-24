package main

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func threadPrompts(threadUrl string) ([]string, error) {
	err := s.Open()
	if err != nil {
		return nil, err
	}

	// Get the thread ID from the URL
	i := strings.LastIndex(threadUrl, "/")
	if len(threadUrl) < i+1 {
		return nil, errors.New("invalid thread URL")
	}
	threadId := threadUrl[i+1:]

	msgs := make([]*discordgo.Message, 0)
	var beforeId string

	for {
		m, err := s.ChannelMessages(threadId, 100, beforeId, "", "")
		if err != nil {
			return "", err
		}

	}

	return "", nil
}
