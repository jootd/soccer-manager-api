package teamdb

import (
	"bytes"
	"strings"

	"github.com/jootd/soccer-manager/business/domain/teambus"
)

func (s *Store) applyFilter(filter teambus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["id"] = filter.ID
		wc = append(wc, "id = :id")
	}

	if filter.Country != nil {
		data["country"] = filter.Country
		wc = append(wc, "country = :country")
	}

	if filter.Name != nil {
		data["name"] = filter.Name
		wc = append(wc, "name = :name")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
