package contracts

import "fmt"

type UsernameExistsError struct {
	Username string
}

func (e *UsernameExistsError) Error() string {
	return fmt.Sprintf("username '%s' already exists", e.Username)
}


type EmailExistsError struct {
	Email string
}

func (e *EmailExistsError) Error() string {
	return fmt.Sprintf("email '%s' already exists", e.Email)
}
