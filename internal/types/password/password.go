package password

import (
	"fmt"
	"regexp"
)

type Password struct {
	value string
}

func (p Password) String() string {
	return p.value
}

var passwordRegEx = regexp.MustCompile("^[a-zA-Z0-9#@!-]{3,19}$")

func Parse(value string) (Password, error) {
	if !passwordRegEx.MatchString(value) {
		return Password{}, fmt.Errorf("invalid password %q", value)
	}

	return Password{
		value: value,
	}, nil
}
