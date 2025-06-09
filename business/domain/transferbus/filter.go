package transferbus

import "github.com/jootd/soccer-manager/business/types/transferstatus"

type QueryFilter struct {
	ID              *int
	PlayerID        *int
	SellerID        *int
	AskingPriceFrom *int64
	AskingPriceTo   *int64
	Status          *transferstatus.Status
}
