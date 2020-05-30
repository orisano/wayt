package main

import (
	"context"
	"flag"
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
		return flag.ErrHelp
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
