package transferdb

import (
	"github.com/jootd/soccer-manager/business/domain/transferbus"
	"github.com/jootd/soccer-manager/business/types/transferstatus"
)

type transfer struct {
	ID          int    `db:"id"`
	PlayerID    int    `db:"player_id"`
	SellerID    int    `db:"seller_id"`
	AskingPrice int64  `db:"asking_price"`
	Status      string `db:"status"`
}

func toBusTransfer(db transfer) transferbus.Transfer {
	return transferbus.Transfer{
		ID:          db.ID,
		PlayerID:    db.PlayerID,
		SellerID:    db.SellerID,
		AskingPrice: db.AskingPrice,
		Status:      transferstatus.MustParse(db.Status),
	}
}

func toDBTransfer(bus transferbus.Transfer) transfer {
	return transfer{
		ID:          bus.ID,
		PlayerID:    bus.PlayerID,
		SellerID:    bus.SellerID,
		AskingPrice: bus.AskingPrice,
		Status:      bus.Status.String(),
	}

}

func toTransferBusSlice(dbTrans []transfer) []transferbus.Transfer {
	var trans []transferbus.Transfer
	for _, v := range dbTrans {
		trans = append(trans, toBusTransfer(v))
	}
	return trans
}
