package teambus

type Team struct {
	ID      int
	Name    string
	Country string
}

type UpdateTeam struct {
	Id      int
	Name    *string
	Country *string
}

type CreateTeam struct {
	Name    string
	Country string
}
