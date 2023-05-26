package buildkitd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

// RootlessDockerProviderOpts represents the options for the buildkitd docker provider.
type RootlessDockerProviderOpts struct {
	Binary string
	Image  string
}

// RootlessDockerProvider implements the buildkit provider via Docker.
type RootlessDockerProvider struct {
	Provider
	Name string
	opts *RootlessDockerProviderOpts
}

// NewRootlessDockerProvider returns a new buildkit provider implemented via Docker.
func NewRootlessDockerProvider(o *RootlessDockerProviderOpts) *RootlessDockerProvider {
	return &RootlessDockerProvider{
		opts: o,
	}
}

// IsSupported implements Provider.IsSupported.
func (p *RootlessDockerProvider) IsSupported(ctx context.Context) error {
	if err := exec.CommandContext(ctx, p.opts.Binary, []string{
		"ps",
	}...).Run(); err != nil {
		return err
	}

	securityOptionsOut, err := exec.CommandContext(ctx, p.opts.Binary, []string{
		"info",
		"--format", `{{ range $opt := .SecurityOptions }}{{ $opt }}{{"\n"}}{{ end }}`,
	}...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not get docker security options: %w\n%s", err, securityOptionsOut)
	}

	if strings.Contains(string(securityOptionsOut), "rootless") {
		return fmt.Errorf("cannot run rootless inside rootless (we're already rootless)")
	}

	driverStatusOut, err := exec.CommandContext(ctx, p.opts.Binary, []string{
		"info",
		"--format", `{{ range $opt := .DriverStatus }}{{ $opt }}{{"\n"}}{{ end }}`,
	}...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not get docker driver status: %w\n%s", err, driverStatusOut)
	}

	if !strings.Contains(string(driverStatusOut), "userxattr true") {
		// I think this is the case. Experiencing this when running Docker in a
		// container on Fedora CoreOS and mounting /var/lib/docker/docker.sock.
		return fmt.Errorf("userxattr=true must be supported by the docker driver")
	}

	return nil
}

// Start implements Provider.Start.
func (p *RootlessDockerProvider) Start(ctx context.Context, address string) error {

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

	runCmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"run",
		"--rm",
		"-d",
		"--security-opt", "seccomp=unconfined",
		"--security-opt", "apparmor=unconfined",
		"--security-opt", "systempaths=unconfined",
		"--name", p.Name,
		"--publish", fmt.Sprintf("%s:%s", portNumber, portNumber),
		p.opts.Image,
		"--addr",
		address,
		"--oci-worker-no-process-sandbox",
	}...)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	log.Info().Str("cmd", strings.Join(runCmd.Args, " ")).Msgf("starting '%s' container", p.Name)
	if err := runCmd.Run(); err != nil {
		return fmt.Errorf("could not run '%s': %w", strings.Join(runCmd.Args, " "), err)
	}
	log.Info().Msgf("started '%s' container", p.Name)

	return nil
}

// Stop implements Provider.Stop.
func (p *RootlessDockerProvider) Stop(ctx context.Context) error {
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

	return nil
}
