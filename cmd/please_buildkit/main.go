package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/urfave/cli/v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		log.Info().Msg("stopping")
		cancel()
	}()

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log_level",
				Value:   "debug",
				EnvVars: []string{"LOG_LEVEL"},
			},
			&cli.StringFlag{
				Name:    "log_format",
				Value:   "auto",
				EnvVars: []string{"LOG_FORMAT"},
			},
		},
		Commands: []*cli.Command{
			WorkerCommand(),
			PushCommand(),
			ReplaceCommand(),
		},
		Before: func(cCtx *cli.Context) error {
			level, err := zerolog.ParseLevel(cCtx.String("log_level"))
			if err != nil {
				return err
			}
			zerolog.SetGlobalLevel(level)
			log.Logger = log.Level(level)

			if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			}

			switch v := cCtx.String("log_format"); {
			case v == "json":
				log.Logger = log.Output(os.Stderr)
			case v == "console":
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			case v == "auto":
				return nil
			default:
				return fmt.Errorf("invalid format: %s", v)
			}

			return nil
		},
	}

	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Fatal().Err(err).Msg("error")
	}
}
