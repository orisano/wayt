package main

import (
	"context"
	"flag"
	"github.com/orisano/subflag"
	"net"
)

type TCPCommand struct {
	address string
}

func (c *TCPCommand) FlagSet() *flag.FlagSet {
	f := flag.NewFlagSet("tcp", flag.ExitOnError)
	f.StringVar(&c.address, "a", c.address, "target address (required)")
	return f
}

func (c *TCPCommand) Run(args []string) error {
	if len(c.address) == 0 {
		return subflag.ErrInvalidArguments
	}
	var d net.Dialer
	for range Loop() {
		conn, err := d.DialContext(context.Background(), "tcp", c.address)
		if err == nil {
			conn.Close()
			return nil
		}
	}
	return ErrTimeout
}
