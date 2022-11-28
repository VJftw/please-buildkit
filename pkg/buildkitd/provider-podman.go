package buildkitd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

// PodmanProviderOpts represents the options for the buildkitd podman provider.
type PodmanProviderOpts struct {
	Binary  string
	Name    string
	Image   string
	Address string
}

// PodmanProvider implements the buildkit provider via Podman.
type PodmanProvider struct {
	Provider
	opts *PodmanProviderOpts
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
func (p *PodmanProvider) Start(ctx context.Context) (string, error) {
	existsCmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"ps",
		"--filter", fmt.Sprintf("name=%s", p.opts.Name),
		"-a",
		"--format", "\"{{.Names}}\"",
	}...)
	out, err := existsCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("could not run '%s': %w", strings.Join(existsCmd.Args, " "), err)
	}

	if strings.Contains(string(out), p.opts.Name) {
		log.Info().Msgf("using existing '%s' container", p.opts.Name)
		return fmt.Sprintf("tcp://%s", p.opts.Address), nil
	}

	log.Info().Msgf("pulling image '%s'", p.opts.Image)
	pullCmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"pull",
		p.opts.Image,
	}...)
	if err := pullCmd.Run(); err != nil {
		return "", fmt.Errorf("could not run '%s': %w", strings.Join(pullCmd.Args, " "), err)
	}

	log.Info().Msgf("starting '%s' container", p.opts.Name)
	portNumber := strings.Split(p.opts.Address, ":")[1]
	runCmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"run",
		"-d",
		"--rm",
		"--name", p.opts.Name,
		"--privileged",
		"--publish", fmt.Sprintf("%s:%s", portNumber, portNumber),
		p.opts.Image,
		"--addr",
		fmt.Sprintf("tcp://%s", p.opts.Address),
	}...)
	if err := runCmd.Run(); err != nil {
		return "", fmt.Errorf("could not run '%s': %w", strings.Join(runCmd.Args, " "), err)
	}
	log.Info().Msgf("started '%s' container", p.opts.Name)

	f, err := os.OpenFile("plz-out/log/please-buildkit-buildkitd.log", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return "", fmt.Errorf("could not open log file")
	}
	if err := f.Truncate(0); err != nil {
		return "", fmt.Errorf("could not truncate log file")
	}
	logsCmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"logs", "-f", p.opts.Name,
	}...)
	logsCmd.Stdout = f
	logsCmd.Stderr = f

	if err := logsCmd.Start(); err != nil {
		return "", fmt.Errorf("could not run '%s': %w", strings.Join(logsCmd.Args, " "), err)
	}

	return fmt.Sprintf("tcp://%s", p.opts.Address), nil
}

// Stop implements Provider.Stop.
func (p *PodmanProvider) Stop(ctx context.Context) error {
	stopCmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"stop",
		p.opts.Name,
	}...)
	if err := stopCmd.Run(); err != nil {
		return fmt.Errorf("could not run '%s': %w", strings.Join(stopCmd.Args, " "), err)
	}

	return nil
}
