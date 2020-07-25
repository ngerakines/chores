package chores

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type loginHandler struct {
	logger  *log.Logger
	persons []string
}

func (h *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(h.logger, w, r, "login", map[string]interface{}{})
		return
	}

	found := false
	name := normalize(r.FormValue("name"))
	for _, person := range h.persons {
		if name == normalize(person) {
			found = true
			break
		}
	}

	if !found {
		renderTemplate(h.logger, w, r, "login", map[string]interface{}{
			"error": "invalid name",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "name", Value: strings.Title(name), Expires: time.Now().AddDate(1, 0, 0)})
	http.Redirect(w, r, "/", http.StatusFound)
}
