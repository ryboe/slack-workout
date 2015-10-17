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

func NewURL(method string, qsp *url.Values) url.URL {
	if qsp == nil {
		qsp = &url.Values{}
	}

	qsp.Set("token", apiToken)
	return url.URL{
		Scheme:   "https",
		Host:     team + ".slack.com",
		Path:     "api/" + method,
		RawQuery: qsp.Encode(),
	}
}

func apiCall(u url.URL, respStruct interface{}) error {
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&respStruct)
	if err != nil {
		return err
	}

	return nil
}

func prettyJSON(js interface{}) (string, error) {
	prettyJs, err := json.MarshalIndent(&js, "", "    ")
	if err != nil {
		return "", err
	}
	return string(prettyJs), nil
}
