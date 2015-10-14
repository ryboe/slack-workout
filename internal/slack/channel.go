// TODO: write package comment
package slack

import (
	"fmt"
	"log"
	"net/url"
	"os"
)

var apiToken string

// TODO: can i get rid of one of these response types?
type channelResponse struct {
	Channel Channel `json:"channel"`
	Ok      bool    `json:"ok"`
	Err     string  `json:"error"`
}

type channelListResponse struct {
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

func init() {
	apiToken = os.Getenv("SLACK_API_TOKEN")
	if apiToken == "" {
		log.Fatal("SLACK_API_TOKEN not set")
	}
}

func NewChannel(team, name string) (Channel, error) {
	var emptyChannel Channel

	qsp := &url.Values{}
	qsp.Set("channel", name)
	listURL := NewURL(team, "channels.list", qsp)
	cl := channelListResponse{}
	err := apiCall(listURL, &cl)
	if err != nil {
		return emptyChannel, err
	}

	if cl.Ok != true {
		return emptyChannel, APIError{cl.Err}
	}

	for _, ch := range cl.Channels {
		if ch.Name == name {
			ch.Team = team
			return ch, nil
		}
	}

	return emptyChannel, fmt.Errorf("no channel named %q on team %q", name, team)
}

func (ch Channel) String() string {
	return fmt.Sprintf("Channel{ID: %s, Name: %s, Members: %v, Team: %s}", ch.ID, ch.Name, ch.Members, ch.Team)
}

func (ch *Channel) UpdateMembers() error {
	qsp := &url.Values{}
	qsp.Set("channel", ch.ID)
	channelURL := NewURL(ch.Team, "channels.info", qsp)

	cr := channelResponse{}
	err := apiCall(channelURL, &cr)
	if err != nil {
		return err
	}

	if !cr.Ok {
		return APIError{cr.Err}
	}

	ch.Members = cr.Channel.Members
	return nil
}
