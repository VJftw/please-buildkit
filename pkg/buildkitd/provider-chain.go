package buildkitd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

// ChainProvider returns a chain implementation of Provider.
type ChainProvider struct {
	providers []Provider
	provider  Provider
}

// NewChainProvider returns a chain provider that implements Provider.
func NewChainProvider(providers ...Provider) *ChainProvider {
	return &ChainProvider{
		providers: providers,
	}
}

// IsSupported implements Provider.IsSupported.
func (p *ChainProvider) IsSupported(ctx context.Context) bool {
	for _, provider := range p.providers {
		if provider.IsSupported(ctx) {
			p.provider = provider
			log.Info().Str("provide", fmt.Sprintf("%T", provider)).Msg("using provider")
			return true
		}
	}

	return false
}

// Start implements Provider.Start.
func (p *ChainProvider) Start(ctx context.Context) (string, error) {
	return p.provider.Start(ctx)
}

// Stop implements Provider.Stop.
func (p *ChainProvider) Stop(ctx context.Context) error {
	return p.provider.Stop(ctx)
}
