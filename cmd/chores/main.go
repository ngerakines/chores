package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"log"
	stdlog "log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/ngerakines/chores"
	"github.com/oklog/run"
)

func main() {
	logger := stdlog.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	c := &chores.Config{}

	flag.StringVar(&c.Address, "address", "0.0.0.0:8080", "The HTTP interface to allow connections on.")
	flag.StringVar(&c.Database, "database", "./chores.db", "The location of the database.")
	flag.BoolVar(&c.Init, "init", false, "Initialize the database.")
	flag.Var(&c.People, "people", "The people that can do chores.")

	flag.Parse()

	if len(c.People) == 0 {
		logger.Fatal("At least one person must be specified.")
	}

	var group run.Group

	db, err := sql.Open("sqlite3", c.Database)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	if c.Init {
		sqlStmt := `CREATE TABLE IF NOT EXISTS events (created_at INTEGER, name TEXT, area TEXT, chore TEXT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			logger.Printf("%q: %s\n", err, sqlStmt)
			return
		}
	}

	ctx := context.Background()

	group.Add(chores.NewServer(logger, db, c))
	group.Add(chores.Run(ctx, logger, chores.NewSignalJob(os.Interrupt, os.Kill)))

	err = group.Run()
	if err != nil {
		if errors.Is(err, chores.SignalError{Signal: os.Interrupt}) {
			return
		} else if errors.Is(err, chores.SignalError{Signal: os.Kill}) {
			return
		}
		logger.Fatal(err)
	}
}
