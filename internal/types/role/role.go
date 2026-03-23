package role

import "fmt"

var (
	Admin = newRole("ADMIN")
	User  = newRole("USER")
)

type Role struct {
	value string
}

var roles = make(map[string]Role)

func newRole(role string) Role {
	r := Role{value: role}

	roles[role] = r

	return r
}

func Parse(value string) (Role, error) {
	role, ok := roles[value]

	if !ok {
		return Role{}, fmt.Errorf("invalid role: %s", value)
	}

	return role, nil
}

func MustParse(value string) Role {
	role, err := Parse(value)

	if err != nil {
		panic(err)
	}

	return role
}

func (r Role) String() string {
	return r.value
}

func ParseToString(r []Role) []string {
	roles := make([]string, len(r))

	for i, v := range r {
		roles[i] = v.String()
	}

	return roles
}

func ParseMany(r []string) ([]Role, error) {
	result := make([]Role, len(r))

	for i, v := range r {
		role, err := Parse(v)

		if err != nil {
			return nil, err
		}

		result[i] = role
	}

	return result, nil
}
