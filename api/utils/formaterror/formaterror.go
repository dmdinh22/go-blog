package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "username") {
		return errors.New("Username has already been used.")
	}

	if strings.Contains(err, "email") {
		return errors.New("Email has already been used.")
	}

	if strings.Contains(err, "title") {
		return errors.New("Title has already been used.")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("The incorrect password has been entered.")
	}

	return errors.New("Incorrect request body.")
}
