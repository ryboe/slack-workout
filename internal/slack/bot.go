// TODO: write package comment
package slack

import (
	"fmt"
	"net/url"
)

type Bot struct {
	Name string
}

func (bot Bot) PostMessage(msg string, ch Channel) error {
	qsp := &url.Values{}
	qsp.Set("channel", ch.ID)
	qsp.Set("icon_url", "https://www.dropbox.com/s/w7hyhjdeylqemx1/sgt_mittens.jpg?dl=0")
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

func (bot Bot) String() string {
	return fmt.Sprintf("%#v", bot)
}

type msgResponse struct {
	Ok  bool   `json:"ok"`
	Err string `json:"error"`
}