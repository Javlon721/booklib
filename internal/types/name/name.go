package name

import (
	"fmt"
	"regexp"
)

type Name struct {
	value string
}

func (n Name) String() string {
	return n.value
}

var nameRegEx = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9' -]{2,19}$")

func Parse(value string) (Name, error) {
	if !nameRegEx.MatchString(value) {
		return Name{}, fmt.Errorf("invalid name %q", value)
	}

	return Name{value: value}, nil
}

func MustParse(value string) Name {
	name, err := Parse(value)

	if err != nil {
		panic(err)
	}

	return name
}
