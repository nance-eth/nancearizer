package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const NANCE_API = "https://api.nance.app/"

type ProposalResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Data    struct {
		Title                 string    `json:"title"`
		Body                  string    `json:"body"`
		Author                string    `json:"author"`
		Coauthors             []string  `json:"coauthors"`
		DiscussionThreadURL   string    `json:"discussionThreadURL"`
		AuthorDiscordId       string    `json:"authorDiscordId"`
		TemperatureCheckVotes []int     `json:"temperatureCheckVotes"`
		CreatedTime           time.Time `json:"createdTime"`
		LastEditedTime        string    `json:"lastEditedTime"`
	} `json:"data"`
}

func proposal(space, proposalId string) (*ProposalResponse, error) {
	resp, err := http.Get(NANCE_API + space + "/proposal/" + proposalId)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	p := &ProposalResponse{}
	if err := json.Unmarshal(bytes, p); err != nil {
		return nil, err
	}

	return p, nil
}
