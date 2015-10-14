// TODO: write package comment
package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const apiURL = "https://%s.slack.com/api/%s?%s"

var apiToken string

// TODO: can i get rid of one of these response types?
type channelResponse struct {
	Channel Channel `json:"channel"`
	Ok      bool    `json:"ok"`
	Err     string  `json:"error,omitempty"`
}

type channelListResponse struct {
	Channels []Channel `json:"channels"`
	Ok       bool      `json:"ok"`
	Err      string    `json:"error,omitempty"`
}

type Channel struct {
	Id      string   `json:"id"`
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

	qsp := map[string]string{
		"channel": name,
		"token":   apiToken,
	}
	listURL := makeURL(apiURL, team, "channels.list", qsp)
	resp, err := http.Get(listURL)
	if err != nil {
		return emptyChannel, err
	}
	defer resp.Body.Close()

	var cl channelListResponse
	err = json.NewDecoder(resp.Body).Decode(&cl)
	if err != nil {
		return emptyChannel, err
	}

	if cl.Ok != true {
		return emptyChannel, APIError{cl.Err}
		// TODO: delete me
		// return emptyChannel, errors.New("failed to get channel list from Slack API")
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
	return fmt.Sprintf("Channel{Id: %s, Name: %s, Members: %v, Team: %s}", ch.Id, ch.Name, ch.Members, ch.Team)
}

func (ch *Channel) UpdateMembers() error {
	qsp := map[string]string{
		"channel": ch.Id,
		"token":   apiToken,
	}
	channelURL := makeURL(apiURL, ch.Team, "channels.info", qsp)

	// TODO: DRY out API calls
	resp, err := http.Get(channelURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	cr := channelResponse{}
	err = json.NewDecoder(resp.Body).Decode(&cr)
	if err != nil {
		return err
	}

	if !cr.Ok {
		return APIError{cr.Err}
	}

	ch.Members = cr.Channel.Members
	return nil
}
