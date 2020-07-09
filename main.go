package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/orisano/subflag"
)

var (
	interval = 1 * time.Second
	timeout  = 5 * time.Minute
)

var (
	deadline time.Time
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("wayt: ")
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flag.DurationVar(&interval, "i", interval, "interval")
	flag.DurationVar(&timeout, "t", timeout, "timeout")
	x := flag.Bool("x", false, "execute command")
	flag.Parse()

	deadline = time.Now().Add(timeout)

	args := flag.Args()
	commands := []subflag.Command{
		&TCPCommand{},
		&SQLCommand{
			driver: "mysql",
			query:  "SELECT 1;",
			envKey: "DB_URL",
		},
		&HTTPCommand{
			method: "GET",
		},
		&FileCommand{},
		&ShellCommand{},
		&GRPCCommand{},
	}

	if len(args) == 0 {
		// for usage
		return subflag.SubCommand(args, commands)
	}
	subCommand := args[0]
	for _, command := range commands {
		flagSet := command.FlagSet()
		if flagSet.Name() != subCommand {
			continue
		}
		if err := flagSet.Parse(args[1:]); err != nil {
			return err
		}

		subArgs := flagSet.Args()
		err := command.Run(nil)
		if err == flag.ErrHelp && flagSet.Usage != nil {
			flagSet.Usage()
		}
		if err != nil {
			return err
		}
		if *x && len(subArgs) > 0 {
			cmd := exec.Command(subArgs[0], subArgs[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Env = os.Environ()
			return cmd.Run()
		}
		return nil
	}
	// for usage
	return subflag.SubCommand(args, commands)
}

func CommandContext() context.Context {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, deadline)
	go func() {
		select {
		case <-sigCh:
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx
}

func Continue(ctx context.Context, interval time.Duration) <-chan struct{} {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	t := time.NewTicker(interval)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-t.C:
			case <-ctx.Done():
				close(ch)
				return
			}
			select {
			case ch <- struct{}{}:
			case <-ctx.Done():
				close(ch)
				return
			}
		}
	}()
	return ch
}
