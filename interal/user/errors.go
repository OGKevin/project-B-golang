package user

import "fmt"

type usernameNotUnique struct {
	username string
}

func (u *usernameNotUnique) Error() string {
	return fmt.Sprintf("Username %s is not unique.", u.username)
}
