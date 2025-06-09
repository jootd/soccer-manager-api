package transferdb

import (
	"bytes"
	"strings"

	"github.com/jootd/soccer-manager/business/domain/transferbus"
)

func (s *Store) applyFilter(filter transferbus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["id"] = filter.ID
		wc = append(wc, "id = :id")
	}

	if filter.PlayerID != nil {
		data["player_id"] = filter.PlayerID
		wc = append(wc, "player_id = :player_id")
	}

	if filter.SellerID != nil {
		data["seller_id"] = filter.SellerID
		wc = append(wc, "seller_id = :seller_id")
	}

	if filter.AskingPriceFrom != nil {
		data["asking_price_from"] = filter.AskingPriceFrom
		wc = append(wc, "asking_price >= :asking_price_from")
	}

	if filter.AskingPriceTo != nil {
		data["asking_price_to"] = filter.AskingPriceTo
		wc = append(wc, "asking_price <= :asking_price_to")
	}

	if filter.Status != nil {
		data["status"] = filter.Status
		wc = append(wc, "status = :status")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
