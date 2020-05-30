package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
)

type HTTPCommand struct {
	urlStr string
	method string
}

func (c *HTTPCommand) FlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("http", flag.ExitOnError)
	flagSet.StringVar(&c.urlStr, "u", c.urlStr, "url (required)")
	flagSet.StringVar(&c.method, "m", c.method, "method")
	return flagSet
}

func (c *HTTPCommand) Run(args []string) error {
	if c.urlStr == "" {
		return flag.ErrHelp
	}
	if _, err := url.ParseRequestURI(c.urlStr); err != nil {
		return fmt.Errorf("parse url: %w", err)
	}
	ctx := CommandContext()
	for range Continue(ctx, interval) {
		permanent, err := c.Request(ctx)
		if permanent {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
	}
	return ctx.Err()
}

func (c *HTTPCommand) Request(ctx context.Context) (bool, error) {
	// TODO(orisano): body to be configurable
	req, err := http.NewRequestWithContext(ctx, c.method, c.urlStr, nil)
	if err != nil {
		return true, fmt.Errorf("new request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		return false, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return true, nil
}
