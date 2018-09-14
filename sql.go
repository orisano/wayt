package main

import (
	"database/sql"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/orisano/subflag"
)

type SQLCommand struct {
	dsn    string
	driver string
	query  string
}

func (c *SQLCommand) FlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("sql", flag.ExitOnError)
	flagSet.StringVar(&c.dsn, "dsn", c.dsn, "data source name (required)")
	flagSet.StringVar(&c.driver, "d", c.driver, "driver")
	flagSet.StringVar(&c.query, "q", c.query, "query")
	return flagSet
}

func (c *SQLCommand) Run(args []string) error {
	if len(c.dsn) == 0 {
		return subflag.ErrInvalidArguments
	}

	for range Loop() {
		db, err := sql.Open(c.driver, c.dsn)
		if err != nil {
			continue
		}
		rows, err := db.Query(c.query)
		if err != nil {
			db.Close()
			continue
		}
		ok := rows.Next()
		rows.Close()
		db.Close()
		if ok {
			return nil
		}
	}
	return ErrTimeout
}
