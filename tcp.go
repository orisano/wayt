package main

import (
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
	ctx := CommandContext()
	var d net.Dialer
	for range Continue(ctx, interval) {
		conn, err := d.DialContext(ctx, "tcp", c.address)
		if err == nil {
			return conn.Close()
		}
		if err := ctx.Err(); err != nil {
			return err
		}
	}
	return ctx.Err()
}
