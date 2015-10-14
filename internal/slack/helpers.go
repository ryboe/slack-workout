// TODO: write package comment
package slack

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type APIError struct {
	msg string
}

func (err APIError) Error() string {
	return fmt.Sprintf("Slack API returned error: %s", err.msg)
}

func makeURL(slackURL, team, method string, qsp map[string]string) string {
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
