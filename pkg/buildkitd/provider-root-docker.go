package buildkitd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

// RootDockerProviderOpts represents the options for the buildkitd docker provider.
type RootDockerProviderOpts struct {
	Binary string
	Image  string
}

// RootDockerProvider implements the buildkit provider via Docker.
type RootDockerProvider struct {
	Provider
	Name string
	opts *RootDockerProviderOpts
}

// NewRootDockerProvider returns a new buildkit provider implemented via Docker.
func NewRootDockerProvider(o *RootDockerProviderOpts) *RootDockerProvider {
	return &RootDockerProvider{
		opts: o,
	}
}

// IsSupported implements Provider.IsSupported.
func (p *RootDockerProvider) IsSupported(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"ps",
	}...)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Info().Err(err).Msg("could not get user home dir, not appending to path")
		cmd.Env = append(
			os.Environ(),
			fmt.Sprintf("PATH=%s/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin", homeDir),
		)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not run '%s': %w\n%s", strings.Join(cmd.Args, " "), err, out)
	}

	return nil
}

// Start implements Provider.Start.
func (p *RootDockerProvider) Start(ctx context.Context, address string) error {

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
	runCmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"run",
		"--rm",
		"-d",
		"--name", p.Name,
		"--privileged",
		"--publish", fmt.Sprintf("%s:%s", portNumber, portNumber),
		p.opts.Image,
		"--addr",
		address,
	}...)

	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	if err := runCmd.Run(); err != nil {
		return fmt.Errorf("could not run '%s': %w", strings.Join(runCmd.Args, " "), err)
	}
	log.Info().Msgf("started '%s' container", p.Name)

	return nil
}

// Stop implements Provider.Stop.
func (p *RootDockerProvider) Stop(ctx context.Context) error {
	log.Info().Msgf("stopping '%s' container", p.Name)
	stopCmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"stop",
		p.Name,
	}...)
	stopCmd.Stdout = os.Stdout
	stopCmd.Stderr = os.Stderr
	if err := stopCmd.Run(); err != nil {
		return fmt.Errorf("could not run '%s': %w", strings.Join(stopCmd.Args, " "), err)
	}

	return nil
}
