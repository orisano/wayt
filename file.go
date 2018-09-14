package main

import (
	"flag"
	"github.com/orisano/subflag"
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
		return subflag.ErrInvalidArguments
	}
	for range Loop() {
		if _, err := os.Lstat(c.path); err == nil {
			return nil
		}
	}
	return ErrTimeout
}
