package main

import (
	"fmt"
	"net"

	"github.com/VJftw/please-buildkit/pkg/buildkitd"
	"github.com/avast/retry-go/v4"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func StartBuildkitdWorker(cCtx *cli.Context) (string, func(), error) {
	chainProvider := buildkitd.NewChainProvider(
		&buildkitd.ChainProviderOpts{},
		buildkitd.NewPodmanProvider(&buildkitd.PodmanProviderOpts{
			Binary: cCtx.String("podman_binary"),
			Image:  cCtx.String("podman_image"),
		}),
		buildkitd.NewRootlessDockerProvider(&buildkitd.RootlessDockerProviderOpts{
			Binary: cCtx.String("rootless_docker_binary"),
			Image:  cCtx.String("rootless_docker_image"),
		}),
		buildkitd.NewRootDockerProvider(&buildkitd.RootDockerProviderOpts{
			Binary: cCtx.String("docker_binary"),
			Image:  cCtx.String("docker_image"),
		}),
	)

	if err := chainProvider.IsSupported(cCtx.Context); err != nil {
		return "", nil, fmt.Errorf("no supported buildkitd providers: %w", err)
	}

	address := ""
	if err := retry.Do(func() error {
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return fmt.Errorf("could not get an available TCP port")
		}

		address = fmt.Sprintf("tcp://0.0.0.0:%d", l.Addr().(*net.TCPAddr).Port)
		if err := l.Close(); err != nil {
			return err
		}

		if err := chainProvider.Start(cCtx.Context, address); err != nil {
			return fmt.Errorf("could not start buildkitd provider: %w", err)
		}

		return nil
	},
		retry.Attempts(10),
		retry.DelayType(retry.BackOffDelay),
		retry.Context(cCtx.Context),
		retry.OnRetry(func(n uint, err error) {
			log.Warn().Msgf("retrying buildkitd worker start")
		})); err != nil {
		return "", nil, err
	}

	if err := buildkitd.WaitForBuildKitWorkers(
		cCtx.String("buildctl_binary"),
		address,
		cCtx.Duration("buildkitd_timeout"),
	); err != nil {
		return "", nil, fmt.Errorf("could not wait for buildkitd workers: %w", err)
	}

	return address, func() {
		if err := chainProvider.Stop(cCtx.Context); err != nil {
			log.Error().Err(err).Msgf("could not stop provider")
		}
	}, nil
}
