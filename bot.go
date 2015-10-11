package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/y0ssar1an/slack-pushups/internal/channel"
)

const slackbotURL = "https://%s.slack.com/services/hooks/slackbot?%s"

var botToken = os.Getenv("SLACK_BOT_TOKEN")

func main() {
	if botToken == "" {
		log.Fatal("SLACK_BOT_TOKEN not set")
	}

	ch, err := channel.New("api-test")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ch)

	// nextMember := make(chan string)
	// go randomMember(ch, nextMember)
}

// // DEBUG
// for i := 0; i < 5; i++ {
// 	fmt.Println(<-nextMember)
// }

// v := url.Values{}
// v.Add("token", botToken)
// v.Add("channel", channel)

// body := bytes.NewBufferString("hello world")
// req, err := http.NewRequest("POST", slackbotURL+v.Encode(), body)
// if err != nil {
// 	log.Fatal(err)
// }

// client := &http.Client{}
// resp, err := client.Do(req)
// if err != nil {
// 	log.Fatal(err)
// }
// defer resp.Body.Close()

// TODO: what is the format of the response?

// apiToken := os.Getenv("SLACK_URL_TOKEN")
// resp, err := http.Get(apiUrl + token)
// if err != nil {
// 	log.Fatal(err)
// }
// defer resp.Body.Close()

// var js interface{}
// err = json.NewDecoder(resp.Body).Decode(&js)
// if err != nil {
// 	log.Fatal(err)
// }

// prettyJs, err := json.MarshalIndent(&js, "", "    ")
// if err != nil {
// 	log.Fatal(err)
// }

// fmt.Println(string(prettyJs))

// func getChannel(channel string) (Channel, error) {
// 	var emptyChannel Channel

// 	qsp := map[string]string{
// 		"channel": channel,
// 		"token":   apiToken,
// 	}
// 	listURL := makeURL(apiURL, "channels.list", qsp)
// 	resp, err := http.Get(listURL)
// 	if err != nil {
// 		return emptyChannel, err
// 	}
// 	defer resp.Body.Close()

// 	var cl ChannelListResponse
// 	err = json.NewDecoder(resp.Body).Decode(&cl)
// 	if err != nil {
// 		return emptyChannel, err
// 	}

// 	if cl.Ok != true {
// 		// DEBUG
// 		fmt.Println(cl)

// 		return emptyChannel, errors.New("failed to get channel list from Slack API")
// 	}

// 	for _, ch := range cl.Channels {
// 		if ch.Name == channel {
// 			return ch, nil
// 		}
// 	}

// 	return emptyChannel, fmt.Errorf("no channel with name %q on team %q", channel, team)
// }

func prettyJSON(js interface{}) (string, error) {
	prettyJs, err := json.MarshalIndent(&js, "", "    ")
	if err != nil {
		return "", err
	}
	return string(prettyJs), nil
}

// func makeURL(slackURL, method string, qsp map[string]string) string {
// 	qs := queryString(qsp)
// 	return fmt.Sprintf(apiURL, team, method, qs)
// }

// func queryString(qsp map[string]string) string {
// 	vals := url.Values{}
// 	for k, v := range qsp {
// 		vals.Add(k, v)
// 	}
// 	return vals.Encode()
// }

func randomMember(ch channel.Channel, nextMember chan string) {
	var err error
	for {
		err = ch.UpdateMembers()
		if err != nil {
			log.Println(err)
		}
		i := rand.Intn(len(ch.Members))
		nextMember <- ch.Members[i]
	}
}

// func updateMembers(members []string, channelURL string) ([]string, bool) {
// 	// If this is the first time running, members will be nil
// 	noMembers := members == nil

// 	resp, err := http.Get(channelURL)
// 	if err != nil {
// 		handleAPIError(noMembers, err)
// 		return members, true
// 	}
// 	defer resp.Body.Close()

// 	body := SlackResponse{}
// 	err = json.NewDecoder(resp.Body).Decode(&body)
// 	if err != nil {
// 		handleAPIError(noMembers, err)
// 		return members, true
// 	}

// 	if !body.Ok {
// 		handleAPIError(noMembers, errors.New("Slack API returned error message"))
// 		return members, true
// 	}

// 	return body.Channel.Members, false
// }
