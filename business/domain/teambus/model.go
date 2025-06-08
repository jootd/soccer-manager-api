package teambus

type Team struct {
	ID      int
	Name    string
	Country string
	Budget  int64
}

type UpdateTeam struct {
	ID      int
	Name    *string
	Country *string
	Budget  *int64
}

type CreateTeam struct {
	Name    string
	Country string
	Budget  int64
}
