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
	Binary  string
	Name    string
	Image   string
	Address string
}

// RootlessDockerProvider implements the buildkit provider via Docker.
type RootlessDockerProvider struct {
	Provider
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
func (p *RootlessDockerProvider) Start(ctx context.Context) (string, error) {
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

	portNumber := strings.Split(p.opts.Address, ":")[1]
	runCmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"run",
		"--rm",
		"-d",
		"--security-opt", "seccomp=unconfined",
		"--security-opt", "apparmor=unconfined",
		"--security-opt", "systempaths=unconfined",
		"--name", p.opts.Name,
		"--publish", fmt.Sprintf("%s:%s", portNumber, portNumber),
		p.opts.Image,
		"--addr",
		fmt.Sprintf("tcp://%s", p.opts.Address),
		"--oci-worker-no-process-sandbox",
	}...)
	log.Info().Str("cmd", strings.Join(runCmd.Args, " ")).Msgf("starting '%s' container", p.opts.Name)
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
func (p *RootlessDockerProvider) Stop(ctx context.Context) error {
	log.Info().Msgf("stopping '%s' container", p.opts.Name)
	stopCmd := exec.CommandContext(ctx, p.opts.Binary, []string{
		"stop",
		p.opts.Name,
	}...)
	if err := stopCmd.Run(); err != nil {
		return fmt.Errorf("could not run '%s': %w", strings.Join(stopCmd.Args, " "), err)
	}

	return nil
}
