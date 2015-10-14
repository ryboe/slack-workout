// TODO: write package comment
package slack

import "fmt"

type userListResponse struct {
	Users []User `json:"members"`
	Ok    bool   `json:"ok"`
	Err   string `json:"error"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Team string // set in NewUser
}

func NewUser(team, name string) (User, error) {
	var emptyUser User

	userURL := NewSlackURL(team, "users.list", nil)
	ur := userListResponse{}
	err := apiCall(userURL, &ur)
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

// TODO: write this
// func (u User) Chat(msg, channel Channel) error {
// }
