package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
)

type ShellCommand struct {
	cmd string
}

func (c *ShellCommand) FlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("sh", flag.ExitOnError)
	flagSet.StringVar(&c.cmd, "c", c.cmd, "command (required)")
	return flagSet
}

func (c *ShellCommand) Run(args []string) error {
	if len(c.cmd) == 0 {
		return flag.ErrHelp
	}
	ctx := CommandContext()
	for range Continue(ctx, interval) {
		cmd := exec.CommandContext(ctx, "/bin/sh", "-c", c.cmd)
		cmd.Stderr = ioutil.Discard
		cmd.Stdout = ioutil.Discard
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("start command: %w", err)
		}
		if err := cmd.Wait(); err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if ok && exitErr.Success() {
				return nil
			}
		} else {
			return nil
		}
	}
	return ctx.Err()
}
