package transferstatus

import "fmt"

type Status struct {
	value string
}

// The set of roles that can be used.
var (
	Listed = newStatus("listed")
	Sold   = newStatus("sold")
)

// Set of known roles.
var statuses = make(map[string]Status)

func newStatus(status string) Status {
	r := Status{status}
	statuses[status] = r
	return r
}

func (t Status) String() string {
	return t.value
}

func (s Status) Equal(t Status) bool {
	return s.value == t.value
}

func (r Status) MarshalText() ([]byte, error) {
	return []byte(r.value), nil
}

func Parse(value string) (Status, error) {
	status, exists := statuses[value]
	if !exists {
		return Status{}, fmt.Errorf("invalid status %q", value)
	}

	return status, nil
}

// useful to check db  values meets business constraints .
func MustParse(value string) Status {
	status, err := Parse(value)
	if err != nil {
		panic(err)
	}

	return status
}
