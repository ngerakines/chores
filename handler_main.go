package chores

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	areas = []string{
		"Entry",
		"Living Room",
		"Downstairs Bathroom",
		"Garage",
		"Kitchen",
		"Office",
		"Study",
		"Back Patio",
		"Front Porch",
		"Front Yard",
		"Driveway",
		"Back Yard",
		"Sun Room",
		"Stairwell",
		"Upstairs Hallway",
		"Hallway",
		"Vanessa's Room",
		"Hannah's Room",
		"Storage Room",
		"Upstairs Bathroom",
		"Master Bedroom",
		"Master Bathroom",
		"Neighborhood",
		"Outside",
		"Put Away",
		"Acura",
		"Honda",
	}
	chores = []string{
		"Vacuumed",
		"Cleaned Surfaces",
		"Swiffered",
		"Watered Plants",
		"Pulled Weeds",
		"Gardened",
		"Organized",
		"Cleaned",
		"Walked Mason",
		"Fed Mason",
		"Put Away Laundry",
		"Washed",
		"Laundry",
		"Dishes",
	}
)

type mainHandler struct {
	logger *log.Logger
	db     *sql.DB
}

func (h *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	nameCookie, err := r.Cookie("name")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		h.logger.Println("error parsing name cookie", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filtersUsed := false
	var filters []string
	for _, param := range []string{"limit", "order", "area", "name", "chore"} {
		if value := r.URL.Query().Get(param); len(value) > 0 {
			filters = append(filters, fmt.Sprintf("%s=%s", param, value))
			filtersUsed = true
		}
	}

	rows, err := queryEvents(h.db, filters...)
	if err != nil {
		h.logger.Println("error getting recent events", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	now := time.Now()

	renderTemplate(h.logger, w, r, "home", map[string]interface{}{
		"now":          now,
		"date_min":     time.Date(now.Year(), now.Month(), 1, now.Hour(), now.Hour(), now.Hour(), 0, now.Location()).Format("2006-01-02"),
		"date_max":     now.AddDate(0, 0, 7).Format("2006-01-02"),
		"date_default": now.Format("2006-01-02"),
		"time_value":   now.Format("15:04"),
		"name":         nameCookie.Value,
		"areas":        areas,
		"chores":       chores,
		"rows":         rows,
		"filters_used": filtersUsed,
	})
}
