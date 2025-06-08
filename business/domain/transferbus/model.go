package transferbus

import "github.com/jootd/soccer-manager/business/types/transferstatus"

type Transfer struct {
	ID          int
	PlayerID    int
	SellerID    int
	AskingPrice int64
	Status      transferstatus.Status
}

type PlayerInfo struct {
	ID    int
	Value int64
}

type TeamInfo struct {
	ID     int
	Budget int64
}

type UpdateTransfer struct {
	ID          int
	PlayerID    *int
	SellerID    *int
	AskingPrice *int64
	Status      *transferstatus.Status
}
