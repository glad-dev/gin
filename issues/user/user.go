package user

import "fmt"

type Details struct {
	Name     string
	Username string
}

func (u *Details) String() string {
	if u.Username == "" {
		return u.Name
	}

	return fmt.Sprintf("%s (%s)", u.Name, u.Username)
}
