package main

import (
	"flag"
	"net/http"
	"net/url"

	"github.com/orisano/subflag"
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
	if _, err := url.ParseRequestURI(c.urlStr); err != nil {
		return subflag.ErrInvalidArguments
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
