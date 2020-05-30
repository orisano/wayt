package main

import (
	"flag"
	"os"
)

type FileCommand struct {
	path string
}

func (c *FileCommand) FlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("file", flag.ExitOnError)
	flagSet.StringVar(&c.path, "p", c.path, "path (required)")
	return flagSet
}

func (c *FileCommand) Run(args []string) error {
	if len(c.path) == 0 {
		return flag.ErrHelp
	}
	ctx := CommandContext()
	for range Continue(ctx, interval) {
		if _, err := os.Lstat(c.path); err == nil {
			return nil
		}
	}
	return ctx.Err()
}
