package buildkitd

import (
	"context"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func WaitForIt(addr string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_ = cancel
	log.Info().Msgf("waiting for %s", addr)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(250 * time.Millisecond):
			if _, err := net.Dial("tcp", addr); err != nil {
				continue
			}
			log.Info().Msgf("%s is available", addr)
			return nil
		}
	}
}

func WaitForGRPC(addr string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_ = cancel
	log.Info().Msgf("waiting for %s", addr)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(250 * time.Millisecond):
			if _, err := net.Dial("tcp", addr); err != nil {
				continue
			}
			if _, err := grpc.Dial(
				addr,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			); err != nil {
				continue
			}
			log.Info().Msgf("%s is available", addr)
			return nil
		}
	}
}
