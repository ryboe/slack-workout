package slack

import (
	"fmt"
	"net/url"
)

type userResponse struct {
	User User   `json:"user"`
	Ok   bool   `json:"ok"`
	Err  string `json:"error"`
}

// A User contains information about a Slack user. It is populated from the
// Slack API.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Team string // set in NewUser
}

// NewUser takes the given Slack user ID (e.g. "U023BECGF") and populates a new
// User from the Slack API.
func NewUser(id string) (User, error) {
	var emptyUser User

	qsp := &url.Values{}
	qsp.Set("user", id)
	userURL := NewURL("users.info", qsp)

	ur := userResponse{}
	err := apiCall(userURL, &ur)
	if err != nil {
		return emptyUser, APIError{err.Error()}
	}

	if !ur.Ok {
		return emptyUser, APIError{ur.Err}
	}

	return ur.User, nil
}

// String returns a human-readable string representation of a User.
func (u User) String() string {
	return fmt.Sprintf("%#v", u)
}
