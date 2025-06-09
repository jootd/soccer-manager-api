package teambus

type Team struct {
	ID      int
	Name    string `json:"name"`
	Country string `json:"country"`
	Budget  int64  `json:"budget"`
}

type UpdateTeam struct {
	ID      int     `json:"id"`
	Name    *string `json:"name"`
	Country *string `json:"country"`
	Budget  *int64  `json:"-"`
}

type CreateTeam struct {
	Name    string
	Country string
	Budget  int64
}
