// TODO: write package comment
package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	team   = "monkeytacos"
	apiURL = "https://%s.slack.com/api/%s?%s"
)

var apiToken string

type channelListResponse struct {
	Channels []Channel `json:"channels"`
	Ok       bool      `json:"ok"`
	Err      string    `json:"error,omitempty"`
}

type channelResponse struct {
	Ok      bool    `json:"ok"`
	Channel Channel `json:"channel"`
	Err     string  `json:"error,omitempty"`
}

type Channel struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

func init() {
	apiToken = os.Getenv("SLACK_API_TOKEN")
	if apiToken == "" {
		log.Fatal("SLACK_API_TOKEN not set")
	}
}

func NewChannel(name string) (Channel, error) {
	var emptyChannel Channel

	qsp := map[string]string{
		"channel": name,
		"token":   apiToken,
	}
	listURL := makeURL(apiURL, "channels.list", qsp)
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
		return emptyChannel, errors.New("failed to get channel list from Slack API")
	}

	for _, ch := range cl.Channels {
		if ch.Name == name {
			return ch, nil
		}
	}

	return emptyChannel, fmt.Errorf("no channel with name %q on team %q", name, team)
}

func (ch Channel) String() string {
	return fmt.Sprintf("Channel{Id: %s, Name: %s, Members: %v}", ch.Id, ch.Name, ch.Members)
}

// TODO: should this be pointer or value?
func (ch *Channel) UpdateMembers() error {
	qsp := map[string]string{
		"channel": ch.Id,
		"token":   apiToken,
	}
	channelURL := makeURL(apiURL, "channels.info", qsp)

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
		return fmt.Errorf("Slack API returned error: %s", cr.Err)
	}

	ch.Members = cr.Channel.Members
	return nil
}

func makeURL(slackURL, method string, qsp map[string]string) string {
	qs := queryString(qsp)
	return fmt.Sprintf(apiURL, team, method, qs)
}

func queryString(qsp map[string]string) string {
	vals := url.Values{}
	for k, v := range qsp {
		vals.Add(k, v)
	}
	return vals.Encode()
}

func prettyJSON(js interface{}) (string, error) {
	prettyJs, err := json.MarshalIndent(&js, "", "    ")
	if err != nil {
		return "", err
	}
	return string(prettyJs), nil
}

type User struct {
	Id   string
	Name string
}

func NewUser(id string) (User, error) {
	// TODO: call slack API to get user name from ID
}
