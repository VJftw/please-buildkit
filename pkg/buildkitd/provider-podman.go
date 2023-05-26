package buildkitd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// PodmanProviderOpts represents the options for the buildkitd podman provider.
type PodmanProviderOpts struct {
	Binary string
	Image  string
}

// PodmanProvider implements the buildkit provider via Podman.
type PodmanProvider struct {
	Provider
	opts *PodmanProviderOpts

	Name string
}

// NewPodmanProvider returns a new buildkit provider implemented via Podman.
func NewPodmanProvider(o *PodmanProviderOpts) *PodmanProvider {
	return &PodmanProvider{
		opts: o,
	}
}

// IsSupported implements Provider.IsSupported.
func (p *PodmanProvider) IsSupported(ctx context.Context) error {
	if err := exec.CommandContext(ctx, p.opts.Binary, []string{
		"ps",
	}...).Run(); err != nil {
		return err
	}

	return nil
}

// Start implements Provider.Start.
func (p *PodmanProvider) Start(ctx context.Context, address string) error {

	portNumber := strings.Split(address, ":")[2]
	p.Name = fmt.Sprintf("please-buildkit-%s", portNumber)

	log.Info().Msgf("pulling image '%s'", p.opts.Image)
	pullCmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"pull",
		p.opts.Image,
	}...)
	if err := pullCmd.Run(); err != nil {
		return fmt.Errorf("could not run '%s': %w", strings.Join(pullCmd.Args, " "), err)
	}

	log.Info().Msgf("starting '%s' container", p.Name)
	// TODO: attempt to set XDG_RUNTIME_DIR, $TMPDIR, $HOME to be much shorter
	runCmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"run",
		"-d",
		"--rm",
		"--name", p.Name,
		"--privileged",
		"--publish", fmt.Sprintf("%s:%s", portNumber, portNumber),
		p.opts.Image,
		"--addr",
		address,
	}...)
	runOut, err := runCmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Strs("cmd", runCmd.Args).Msgf("%s", runOut)
		return err
	}
	log.Info().Msgf("started '%s' container", p.Name)

	return nil
}

// Stop implements Provider.Stop.
func (p *PodmanProvider) Stop(ctx context.Context) error {
	log.Info().Msgf("stopping '%s' container", p.Name)
	stopCmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"stop",
		p.Name,
	}...)
	stopOut, err := stopCmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Strs("cmd", stopCmd.Args).Msgf("%s", stopOut)
		return err
	}

	cleanupCmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"unshare",
		"rm", "-rf", filepath.Join(os.Getenv("HOME"), ".local/share/containers/storage"),
	}...)
	cleanupOut, err := cleanupCmd.CombinedOutput()
	if err != nil {
		log.Warn().Err(err).Strs("cmd", cleanupCmd.Args).Msgf("%s", cleanupOut)
	}

	return nil
}
