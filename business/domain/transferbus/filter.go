package transferbus

import "github.com/jootd/soccer-manager/business/types/transferstatus"

type QueryFilter struct {
	ID          *int
	PlayerID    *int
	SellerID    *int
	AskingPrice *int64
	Status      *transferstatus.Status
}
