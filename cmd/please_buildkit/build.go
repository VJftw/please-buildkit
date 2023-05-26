package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func BuildCommand() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Builds a Docker Image directly",
		Description: `
This command builds a docker image directly with the given parameters as Please
> 17.0.0 does not support Please workers anymore.
`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "image_out",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "fqn_tags_file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dockerfile",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "buildctl_binary",
				Value: "buildctl",
			},
			&cli.DurationFlag{
				Name:  "buildkitd_timeout",
				Value: 5 * time.Second,
			},
			&cli.StringFlag{
				Name:  "docker_binary",
				Value: "docker",
			},
			&cli.StringFlag{
				Name:  "docker_image",
				Value: "moby/buildkit:master",
			},
			&cli.StringFlag{
				Name:  "rootless_docker_binary",
				Value: "docker",
			},
			&cli.StringFlag{
				Name:  "rootless_docker_image",
				Value: "moby/buildkit:master-rootless",
			},
			&cli.StringFlag{
				Name:  "podman_binary",
				Value: "podman",
			},
			&cli.StringFlag{
				Name:  "podman_image",
				Value: "docker.io/moby/buildkit:master",
			},
		},
		Action: func(cCtx *cli.Context) error {

			buildkitdAddr, closeFn, err := StartBuildkitdWorker(cCtx)
			if err != nil {
				return err
			}
			defer closeFn()

			tmpDir := os.TempDir()

			fqnTagsFileContents, err := os.ReadFile(cCtx.String("fqn_tags_file"))
			if err != nil {
				return fmt.Errorf("could not read '%s': %w", cCtx.String("fqn_tags_file"), err)
			}
			fqnTags := strings.FieldsFunc(string(fqnTagsFileContents), func(c rune) bool {
				return c == '\n'
			})

			if err := os.MkdirAll(filepath.Join(tmpDir, "dockerfile"), 0755); err != nil {
				return fmt.Errorf("could not create 'dockerfile' dir: %w", err)
			}

			if err := os.Rename(
				filepath.Join(tmpDir, cCtx.String("dockerfile")),
				filepath.Join(tmpDir, "dockerfile/Dockerfile"),
			); err != nil {
				return fmt.Errorf("could not move dockerfile: %w", err)
			}

			outImagePath := cCtx.String("image_out")
			cmd := exec.CommandContext(cCtx.Context,
				cCtx.String("buildctl_binary"),
				[]string{
					"build",
					"--frontend=dockerfile.v0",
					"--no-cache",
					"--trace", filepath.Join(tmpDir, "buildctl.trace"),
					"--local", fmt.Sprintf("context=%s", tmpDir),
					"--local", fmt.Sprintf("dockerfile=%s", filepath.Join(tmpDir, "dockerfile")),
					"--output", fmt.Sprintf("type=docker,\"name=%s\",dest=%s", strings.Join(fqnTags, ","), outImagePath),
				}...)

			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			cmd.Env = append(os.Environ(), []string{
				"BUILDKIT_HOST=" + buildkitdAddr,
			}...)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("could not run '%s': %w\n%s", strings.Join(cmd.Args, " "), err, stderr.String())
			}

			log.Info().
				Str("out", outImagePath).
				Msg("built image")

			return nil
		},
	}
}
