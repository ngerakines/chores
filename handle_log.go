package chores

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

type logHandler struct {
	logger *log.Logger
	db     *sql.DB
}

func (h *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	nameCookie, err := r.Cookie("name")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.logger.Printf("unable to parse form: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	area := r.FormValue("area")
	chore := r.FormValue("chore")
	when, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", r.FormValue("date"), r.FormValue("time")), time.Local)
	if err != nil {
		h.logger.Printf("unable to parse date/time %s %s: %s", r.FormValue("date"), r.FormValue("time"), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Printf("storing event %s", when)

	sqlStmt := `INSERT INTO events(created_at, name, area, chore) VALUES (?, ?, ?, ?)`
	_, err = h.db.Exec(sqlStmt, when.Unix(), nameCookie.Value, area, chore)
	if err != nil {
		h.logger.Printf("unable to store event: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
