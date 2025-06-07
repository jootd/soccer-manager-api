package playerdb

import (
	"bytes"
	"strings"

	"github.com/jootd/soccer-manager/business/domain/playerbus"
)

func (s *Store) applyFilter(filter playerbus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["id"] = filter.ID
		wc = append(wc, "id = :id")
	}

	if filter.TeamId != nil {
		data["team_id"] = filter.TeamId
		wc = append(wc, "team_id = :team_id")
	}

	if filter.FirstName != nil {
		data["name"] = filter.FirstName
		wc = append(wc, "name = :first_name")
	}

	if filter.LastName != nil {
		data["last_name"] = filter.LastName
		wc = append(wc, "last_name = :last_name")
	}

	if filter.ValueFrom != nil {
		data["value_from"] = filter.ValueFrom
		wc = append(wc, "value >= :value_from")
	}

	if filter.ValueTo != nil {
		data["value_to"] = filter.ValueTo
		wc = append(wc, "value <= :value_to")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
