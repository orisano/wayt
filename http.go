package main

import (
	"flag"
	"net/http"
	"net/url"
)

type HTTPCommand struct {
	urlStr string
}

func (c *HTTPCommand) FlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("http", flag.ExitOnError)
	flagSet.StringVar(&c.urlStr, "u", c.urlStr, "url (required)")
	return flagSet
}

func (c *HTTPCommand) Run(args []string) error {
	if c.urlStr == "" {
		return flag.ErrHelp
	}
	if _, err := url.ParseRequestURI(c.urlStr); err != nil {
		return fmt.Errorf("parse url: %w", err)
	}
	for range Loop() {
		resp, err := http.Get(c.urlStr)
		if err != nil {
			continue
		}
		resp.Body.Close()
		if resp.StatusCode < 500 {
			return nil
		}
	}
	return ErrTimeout
}
