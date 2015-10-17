// Package slack implements the Channel type, for getting information about
// Slack channels, and the User type, for sending chats.
package slack

import (
	"fmt"
	"log"
	"net/url"
	"os"
)

const team = "monkeytacos"

var apiToken string

func init() {
	apiToken = os.Getenv("SLACK_API_TOKEN")
	if apiToken == "" {
		log.Fatal("SLACK_API_TOKEN not set")
	}
}

type channelResponse struct {
	Channel  Channel   `json:"channel"`
	Channels []Channel `json:"channels"`
	Ok       bool      `json:"ok"`
	Err      string    `json:"error"`
}

type Channel struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
	Team    string   // set in NewChannel
}

func NewChannel(name string) (Channel, error) {
	var emptyChannel Channel

	qsp := &url.Values{}
	qsp.Set("channel", name)
	listURL := NewURL("channels.list", qsp)
	cr := channelResponse{}
	err := apiCall(listURL, &cr)
	if err != nil {
		return emptyChannel, APIError{err.Error()}
	}

	if cr.Ok != true {
		return emptyChannel, APIError{cr.Err}
	}

	for _, ch := range cr.Channels {
		if ch.Name == name {
			ch.Team = team
			return ch, nil
		}
	}

	return emptyChannel, fmt.Errorf("no channel named %q on team %q", name, team)
}

func (ch Channel) String() string {
	return fmt.Sprintf("%#v", ch)
}

func (ch *Channel) UpdateMembers() error {
	qsp := &url.Values{}
	qsp.Set("channel", ch.ID)
	channelURL := NewURL("channels.info", qsp)

	cr := channelResponse{}
	err := apiCall(channelURL, &cr)
	if err != nil {
		return APIError{err.Error()}
	}

	if !cr.Ok {
		return APIError{cr.Err}
	}

	ch.Members = cr.Channel.Members
	return nil
}
