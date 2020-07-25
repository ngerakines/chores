package chores

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// NewServer is neat.
func NewServer(logger *log.Logger, db *sql.DB, config *Config) (func() error, func(error)) {

	fs := http.FileServer(AssetFile())

	mux := http.NewServeMux()

	mux.Handle("/static/", fs)
	mux.Handle("/login", &loginHandler{persons: config.People, logger: logger})
	mux.Handle("/log", &logHandler{logger: logger, db: db})
	mux.Handle("/", &mainHandler{logger: logger, db: db})

	srv := http.Server{
		Addr:    config.Address,
		Handler: loggingMiddleware(logger)(mux),
	}

	return func() error {
			logger.Println("starting http service", srv.Addr)
			return srv.ListenAndServe()
		}, func(error) {
			httpCtx, httpCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer httpCancel()
			logger.Println("stopping http service")
			if err := srv.Shutdown(httpCtx); err != nil {
				logger.Println("error stopping http service", err)
			}
		}
}

func renderTemplate(logger *log.Logger, w http.ResponseWriter, r *http.Request, target string, data map[string]interface{}) {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", fmt.Sprintf("%s.html", target))

	info, err := AssetInfo(fp)
	if err != nil {
		logger.Println("unable to find template", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		logger.Println("target is a directory")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	funcs := map[string]interface{}{
		"datetime": func(in time.Time) string {
			return in.Local().Format("02 Jan 2006 15:04")
		},
	}

	tmpl, err := NewTemplate("", Asset).Funcs(funcs).ParseFiles(lp, fp)
	if err != nil {
		logger.Println("error parsing template", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		logger.Println("error parsing template", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
