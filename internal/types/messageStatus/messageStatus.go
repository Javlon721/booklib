package messageStatus

import "fmt"

var (
	Pending  = newStatus("PENDING")
	Recieved = newStatus("RECIEVED")
	Read     = newStatus("READ")
)

type Status struct {
	value string
}

var statuses = make(map[string]Status)

func newStatus(status string) Status {
	r := Status{value: status}

	statuses[status] = r

	return r
}

func Parse(value string) (Status, error) {
	status, ok := statuses[value]

	if !ok {
		return Status{}, fmt.Errorf("invalid status: %s", value)
	}

	return status, nil
}

func MustParse(value string) Status {
	status, err := Parse(value)

	if err != nil {
		panic(err)
	}

	return status
}

func (r Status) String() string {
	return r.value
}

func ParseToString(s []Status) []string {
	result := make([]string, len(s))

	for i, v := range s {
		result[i] = v.String()
	}

	return result
}

func ParseMany(s []string) ([]Status, error) {
	result := make([]Status, len(s))

	for i, v := range s {
		status, err := Parse(v)

		if err != nil {
			return nil, err
		}

		result[i] = status
	}

	return result, nil
}
