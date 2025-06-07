package playerbus

type Player struct {
	ID        int
	TeamID    int
	FirstName string
	LastName  string
	Age       int
	Country   string
	Value     int64
}

type UpdatePlayer struct {
	FirstName *string
	LastName  *string
	Country   *string
	Value     *int64
}
