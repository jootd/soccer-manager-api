package playerbus

type Player struct {
	ID        int
	FirstName string
	LastName  string
	Country   string
	Value     float64
	TeamId    int
}

type UpdatePlayer struct {
	FirstName *string
	LastName  *string
	Country   *string
}
