package playerbus

type QueryFilter struct {
	ID        *int
	TeamId    *int
	FirstName *string
	LastName  *string
	Country   *string
	ValueFrom *int64
	ValueTo   *int64
}
