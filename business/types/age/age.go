package age

import "fmt"

type Age struct {
	value int
}

func (a Age) Value() int {
	return a.value
}

func (a Age) String() string {
	return fmt.Sprintf("%d", a.value)
}

func (m Age) Equal(a2 Age) bool {
	return m.value == a2.value
}

func (a Age) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func Parse(value int) (Age, error) {
	if value < 18 || value > 40 {
		return Age{}, fmt.Errorf("invalid age %d", value)
	}

	return Age{value}, nil
}

func MustParse(value int) Age {
	age, err := Parse(value)
	if err != nil {
		panic(err)
	}

	return age
}
