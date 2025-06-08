package playeradapter

import (
	"github.com/jootd/soccer-manager/business/domain/playerbus"
	"github.com/jootd/soccer-manager/business/domain/transferbus"
)

func toTransferInfo(p playerbus.Player) transferbus.PlayerInfo {
	return transferbus.PlayerInfo{
		ID:    p.ID,
		Value: p.Value,
	}
}
