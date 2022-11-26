package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/VJftw/please-buildkit/internal/cmd"
	"github.com/VJftw/please-buildkit/pkg/buildkitd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func WorkerCommand() *cli.Command {
	return &cli.Command{
		Name:  "worker",
		Usage: "Starts as a Please worker for building Docker Images via BuildKit",
		Description: `
This command starts the Please worker for building docker images via Buildkit.
Please note there is no stdout from this command as Please registers stdout from
workers as errors. Logs from this command are written to 'plz-out/log/'.
`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "buildctl_binary",
				Value: "buildctl",
			},
			&cli.StringFlag{
				Name:  "docker_binary",
				Value: "docker",
			},
			&cli.StringFlag{
				Name:  "docker_name",
				Value: "please-buildkitd",
			},
			&cli.StringFlag{
				Name:  "docker_image",
				Value: "moby/buildkit:master",
			},
			&cli.StringFlag{
				Name:  "podman_binary",
				Value: "podman",
			},
			&cli.StringFlag{
				Name:  "podman_name",
				Value: "please-buildkitd",
			},
			&cli.StringFlag{
				Name:  "podman_image",
				Value: "docker.io/moby/buildkit:master",
			},
			&cli.StringFlag{
				Name:  "buildkitd_address",
				Value: "0.0.0.0:1234",
			},
		},
		Action: func(cCtx *cli.Context) error {
			f, err := os.OpenFile("plz-out/log/please-buildkit-worker.log", os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				log.Fatal().Err(err).Msg("could not open log file")
			}
			defer f.Close()
			if err := f.Truncate(0); err != nil {
				log.Fatal().Err(err).Msg("could not truncate log file")
			}

			multi := zerolog.MultiLevelWriter(
				cmd.NewLevelWriter(zerolog.FatalLevel, log.Logger.With().Logger()),
				cmd.NewLevelWriter(zerolog.InfoLevel, zerolog.New(f).With().Logger()),
			)
			log.Logger = zerolog.New(multi).With().Timestamp().Logger()

			chainProvider := buildkitd.NewChainProvider(
				buildkitd.NewDockerProvider(&buildkitd.DockerProviderOpts{
					Binary:  cCtx.String("docker_binary"),
					Name:    cCtx.String("docker_name"),
					Image:   cCtx.String("docker_image"),
					Address: cCtx.String("buildkitd_address"),
				}),
				buildkitd.NewPodmanProvider(&buildkitd.PodmanProviderOpts{
					Binary:  cCtx.String("podman_binary"),
					Name:    cCtx.String("podman_name"),
					Image:   cCtx.String("podman_image"),
					Address: cCtx.String("buildkitd_address"),
				}),
			)

			if !chainProvider.IsSupported(cCtx.Context) {
				return fmt.Errorf("no supported buildkitd providers")
			}

			buildkitdAddr, err := chainProvider.Start(cCtx.Context)
			if err != nil {
				return fmt.Errorf("could not start buildkitd provider: %w", err)
			}

			pleaseWorker := buildkitd.NewPleaseWorker(&buildkitd.PleaseWorkerOpts{
				BuildKitAddress: buildkitdAddr,
				BuildCtlBinary:  cCtx.String("buildctl_binary"),
			})

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := pleaseWorker.Start(cCtx.Context); err != nil {
					log.Error().Err(err).Msg("could not start please worker")
				}
				if err := chainProvider.Stop(context.Background()); err != nil {
					log.Error().Err(err).Msg("could not stop provider")
				}
			}()
			log.Info().Msg("started please worker")
			wg.Wait()
			log.Info().Msg("please worker has stopped")

			return nil
		},
	}
}
