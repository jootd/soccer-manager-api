package position

import "fmt"

// The set of roles that can be used.
var (
	Goalkeeper = newPosition("Goalkeeper")
	Defender   = newPosition("Defender")
	Midfielder = newPosition("Midfielder")
	Attacker   = newPosition("Attacker")
)

// Set of known roles.
var positions = make(map[string]Position)

type Position struct {
	value string
}

func newPosition(position string) Position {
	r := Position{position}
	positions[position] = r
	return r
}

func (p Position) String() string {
	return p.value
}

func (r Position) Equal(r2 Position) bool {
	return r.value == r2.value
}

func (r Position) MarshalText() ([]byte, error) {
	return []byte(r.value), nil
}

func Parse(value string) (Position, error) {
	position, exists := positions[value]
	if !exists {
		return Position{}, fmt.Errorf("invalid role %q", value)
	}

	return position, nil
}

// useful to check db  values meets business constraints .
func MustParse(value string) Position {
	role, err := Parse(value)
	if err != nil {
		panic(err)
	}

	return role
}
