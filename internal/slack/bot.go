package slack

import (
	"fmt"
	"net/url"
)

// A Bot is a fakey Slack user that can post messages.
type Bot struct {
	Name string
}

// PostMessage posts the given message to the given Channel.
func (bot *Bot) PostMessage(msg string, ch Channel) error {
	qsp := &url.Values{}
	qsp.Set("channel", "#"+ch.Name)
	qsp.Set("icon_url", "http://i.imgur.com/XpotArk.jpg")
	qsp.Set("link_names", "1")
	qsp.Set("text", msg)
	qsp.Set("username", bot.Name)
	botURL := NewURL("chat.postMessage", qsp)

	mr := msgResponse{}
	err := apiCall(botURL, &mr)
	if err != nil {
		return APIError{err.Error()}
	}

	if !mr.Ok {
		return APIError{mr.Err}
	}

	return nil
}

// String returns a human-readable string representation of a Bot.
func (bot Bot) String() string {
	return fmt.Sprintf("%#v", bot)
}

type msgResponse struct {
	Ok  bool   `json:"ok"`
	Err string `json:"error"`
}
