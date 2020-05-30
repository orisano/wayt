package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/xo/dburl"
)

type SQLCommand struct {
	dsn    string
	driver string
	query  string
	url    string
	envKey string
}

func (c *SQLCommand) FlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("sql", flag.ExitOnError)
	flagSet.StringVar(&c.dsn, "dsn", c.dsn, "data source name (required)")
	flagSet.StringVar(&c.driver, "d", c.driver, "driver")
	flagSet.StringVar(&c.query, "q", c.query, "query")
	flagSet.StringVar(&c.url, "url", c.url, "url")
	flagSet.StringVar(&c.envKey, "env", c.envKey, "")
	return flagSet
}

func (c *SQLCommand) Run(args []string) error {
	if c.url == "" {
		c.url = os.Getenv(c.envKey)
	}
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

	ctx := CommandContext()
	for range Continue(ctx, interval) {
		permanent, err := c.Query(ctx)
		if permanent {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
	}
	return ctx.Err()
}

func (c *SQLCommand) Query(ctx context.Context) (bool, error) {
	db, err := sql.Open(c.driver, c.dsn)
	if err != nil {
		return true, fmt.Errorf("open db: %w", err)
	}
	defer db.Close()
	rows, err := db.QueryContext(ctx, c.query)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return false, rows.Err()
	}
	return true, nil
}
