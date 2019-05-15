package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/orisano/subflag"
)

var (
	interval = 1 * time.Second
	timeout  = 5 * time.Minute
)

var ErrTimeout = errors.New("timeout")

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
	flag.Parse()
	return subflag.SubCommand(flag.Args(), []subflag.Command{
		&TCPCommand{},
		&SQLCommand{
			driver: "mysql",
			query:  "SELECT 1;",
		},
		&HTTPCommand{},
		&FileCommand{},
		&ShellCommand{},
	})
}

func Loop() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGTERM)

		timer := time.NewTimer(timeout)
		defer timer.Stop()

		ch <- struct{}{}

		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				ch <- struct{}{}
			case <-timer.C:
				close(ch)
				return
			case <-sigCh:
				close(ch)
				return
			}
		}
	}()
	return ch
}
