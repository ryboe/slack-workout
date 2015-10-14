// TODO: write package comment
package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type APIError struct {
	msg string
}

func (err APIError) Error() string {
	return fmt.Sprintf("Slack API returned error: %s", err.msg)
}

type SlackURL struct {
	u url.URL
}

func NewSlackURL(team, method string, qsp *url.Values) SlackURL {
	if qsp == nil {
		qsp = &url.Values{}
	}

	qsp.Set("token", apiToken)
	return SlackURL{
		u: url.URL{
			Scheme:   "https",
			Host:     team + ".slack.com",
			Path:     "api/" + method,
			RawQuery: qsp.Encode(),
		},
	}
}

func apiCall(su SlackURL, respStruct interface{}) error {
	resp, err := http.Get(su.u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&respStruct)
	if err != nil {
		return err
	}

	if !respStruct.Ok {
		return APIError{respStruct.Err}
	}

	return nil
}

// TODO: delete this
// func makeURL(slackURL, team, method string, qsp map[string]string) string {
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

func prettyJSON(js interface{}) (string, error) {
	prettyJs, err := json.MarshalIndent(&js, "", "    ")
	if err != nil {
		return "", err
	}
	return string(prettyJs), nil
}
