package buildkitd

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

type ChainProviderOpts struct {
}

// ChainProvider returns a chain implementation of Provider.
type ChainProvider struct {
	opts      *ChainProviderOpts
	providers []Provider
	provider  Provider
}

// NewChainProvider returns a chain provider that implements Provider.
func NewChainProvider(opts *ChainProviderOpts, providers ...Provider) *ChainProvider {
	return &ChainProvider{
		opts:      opts,
		providers: providers,
	}
}

// IsSupported implements Provider.IsSupported.
func (p *ChainProvider) IsSupported(ctx context.Context) error {
	allErrs := []error{}

	for _, provider := range p.providers {
		if err := provider.IsSupported(ctx); err != nil {
			log.Warn().Err(err).Msgf("%T is unsupported", provider)
			allErrs = append(allErrs, err)
		} else {
			p.provider = provider
			log.Info().Str("provider", fmt.Sprintf("%T", provider)).Msg("using provider")
			return nil
		}
	}

	return errors.Join(allErrs...)
}

// Start implements Provider.Start.
func (p *ChainProvider) Start(ctx context.Context, address string) error {
	return p.provider.Start(ctx, address)
}

// Stop implements Provider.Stop.
func (p *ChainProvider) Stop(ctx context.Context) error {
	return p.provider.Stop(ctx)
}
