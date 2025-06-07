package playerbus

type QueryFilter struct {
	ID        *int
	TeamId    *int
	FirstName *string
	LastName  *string
	Country   *string
	ValueFrom *float64
	ValueTo   *float64
}
