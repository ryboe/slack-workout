// TODO: write package comment
package slack

import "fmt"

type userResponse struct {
	User User   `json:"user"`
	Ok   bool   `json:"ok"`
	Err  string `json:"error"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Team string // set in NewUser
}

func NewUser(id string) (User, error) {
	var emptyUser User

	userURL := NewURL("users.info", nil)
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

func (u User) String() string {
	return fmt.Sprintf("%#v", u)
}
