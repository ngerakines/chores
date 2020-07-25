package chores

import (
	"database/sql"
	"strings"
	"time"
)

// Event is neat.
type Event struct {
	createdAt int64
	Who       string
	Area      string
	Chore     string
}

// When is neat.
func (e *Event) When() time.Time {
	return time.Unix(e.createdAt, 0)
}

func validateUser(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)
	switch input {
	case "nick":
		return "Nick"
	case "mattie":
		return "Mattie"
	case "vanessa":
		return "Vanessa"
	case "hannah":
		return "Hannah"
	default:
		return ""
	}
}

func queryEvents(db *sql.DB, filters ...string) ([]*Event, error) {
	var results []*Event

	var query strings.Builder

	var where []string
	var whereValues []interface{}
	order := ""
	limit := ""

	query.WriteString("SELECT created_at, name, area, chore FROM events")
	for _, filter := range filters {
		if strings.HasPrefix(filter, "limit=") {
			limit = strings.TrimPrefix(filter, "limit=")

		} else if strings.HasPrefix(filter, "order=") {
			order = strings.TrimPrefix(filter, "order=")

		} else if strings.HasPrefix(filter, "area=") {
			value := strings.TrimPrefix(filter, "area=")
			where = append(where, "area = ?")
			whereValues = append(whereValues, value)

		} else if strings.HasPrefix(filter, "name=") {
			value := strings.TrimPrefix(filter, "name=")
			where = append(where, "name = ?")
			whereValues = append(whereValues, value)

		} else if strings.HasPrefix(filter, "chore=") {
			value := strings.TrimPrefix(filter, "chore=")
			where = append(where, "chore = ?")
			whereValues = append(whereValues, value)
		}
	}
	if len(filters) == 0 {
		order = "created_at"
		limit = "50"
	}

	if len(where) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(where, ", "))
	}

	if len(order) > 0 {
		query.WriteString(" ORDER BY ")
		query.WriteString(order)
	}
	if len(limit) > 0 {
		query.WriteString(" LIMIT ")
		query.WriteString(limit)
	}

	rows, err := db.Query(query.String(), whereValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		e := &Event{}
		err = rows.Scan(&e.createdAt, &e.Who, &e.Area, &e.Chore)
		if err != nil {
			return nil, err
		}
		results = append(results, e)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}
