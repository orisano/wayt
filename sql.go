package main

import (
	"database/sql"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/xo/dburl"
)

type SQLCommand struct {
	dsn    string
	driver string
	query  string
	url    string
}

func (c *SQLCommand) FlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("sql", flag.ExitOnError)
	flagSet.StringVar(&c.dsn, "dsn", c.dsn, "data source name (required)")
	flagSet.StringVar(&c.driver, "d", c.driver, "driver")
	flagSet.StringVar(&c.query, "q", c.query, "query")
	flagSet.StringVar(&c.url, "url", c.url, "url")
	return flagSet
}

func (c *SQLCommand) Run(args []string) error {
	if c.dsn == "" && c.url == "" {
		return flag.ErrHelp
	}
	if c.url != "" {
		u, err := dburl.Parse(c.url)
		if err != nil {
			return fmt.Errorf("parse dburl: %w", err)
		}
		c.driver = u.Driver
		c.dsn = u.DSN
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
