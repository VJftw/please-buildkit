package buildkitd

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
)

// Provider abstracts the implementations of BuildKitD providers that run
// `buildkitd` as a daemon.
type Provider interface {
	// IsSupported returns whether or not the implementation is supported on
	// this host.
	IsSupported(ctx context.Context) bool
	// Start starts the `buildkitd` daemon using the implementation and returns
	// the buildkitd address to use as `BUILDKIT_HOST`. This should wait for the
	// daemon to be ready.
	Start(ctx context.Context) (string, error)
	// Stop stops the `buildkitd` daemon using the implementation.
	Stop(ctx context.Context) error
}

func WaitForBuildKitWorkers(
	buildctlBinary string,
	addr string,
	timeout time.Duration,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_ = cancel
	log.Info().Msgf("waiting for buildkit workers %s", addr)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(250 * time.Millisecond):
			cmd := exec.CommandContext(ctx,
				buildctlBinary,
				[]string{"debug", "workers"}...,
			)
			cmd.Env = append(os.Environ(), []string{"BUILDKIT_HOST=" + addr}...)
			stdoutStderr, err := cmd.CombinedOutput()
			if err != nil {
				log.Warn().Err(err).Msgf("buildkit workers failed: %s", stdoutStderr)
				continue
			}

			log.Info().Msgf("%s is available", addr)

			return nil
		}
	}
}
