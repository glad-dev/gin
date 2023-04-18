package remote

import "fmt"

type User struct {
	Name     string
	Username string
}

func (u *User) String() string {
	if u.Username == "" {
		return u.Name
	}

	return fmt.Sprintf("%s (%s)", u.Name, u.Username)
}
