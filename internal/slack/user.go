// TODO: write package comment
package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO: can i get rid of one of these types?
type userResponse struct {
	User User   `json:"user"`
	Ok   bool   `json:"ok"`
	Err  string `json:"error,omitempty"`
}

type userListResponse struct {
	Users []User `json:"users"`
	Ok    bool   `json:"ok"`
	Err   string `json:"error,omitempty"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Team string // set in NewUser
}

func NewUser(team, name string) (User, error) {
	var emptyUser User

	qsp := map[string]string{
		"token": apiToken,
	}
	userURL := makeURL(apiURL, team, "users.list", qsp)

	resp, err := http.Get(userURL)
	if err != nil {
		return emptyUser, err
	}
	defer resp.Body.Close()

	ur := userListResponse{}
	err = json.NewDecoder(resp.Body).Decode(&ur)
	if err != nil {
		return emptyUser, err
	}

	if !ur.Ok {
		return emptyUser, APIError{ur.Err}
	}

	for _, u := range ur.Users {
		if u.Name == name {
			u.Team = team
			return u, nil
		}
	}

	return emptyUser, fmt.Errorf("no user named %q on team %q", name, team)
}

func (u User) String() string {
	return fmt.Sprintf("User{Id: %s, Name: %s, Team: %s}", u.Id, u.Name, u.Team)
}
