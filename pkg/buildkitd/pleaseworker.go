package buildkitd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

type PleaseWorkerOpts struct {
	BuildKitAddress string
	BuildCtlBinary  string `long:"buildctl_binary" env:"BUILDCTL_BINARY" description:"the path to the binary to use as 'buildctl'"`
}

type PleaseWorker struct {
	opts *PleaseWorkerOpts
}

func NewPleaseWorker(o *PleaseWorkerOpts) *PleaseWorker {
	return &PleaseWorker{
		opts: o,
	}
}

func (w *PleaseWorker) Start(ctx context.Context) error {
	decodedMsg := make(chan *Request, 1)
	decoder := json.NewDecoder(os.Stdin)
	go func() {
		for {
			select {
			case <-ctx.Done():
			default:
				msg := &Request{}
				if err := decoder.Decode(msg); err != nil {
					log.Error().Err(err).Msg("could not decode stdin")
					return
				}
				decodedMsg <- msg
			}

		}

	}()

	encoder := json.NewEncoder(os.Stdout)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-decodedMsg:
			go func() {
				log.Info().Str("rule", msg.Rule).Msg("handling request")
				respMsg := &Response{
					Rule: msg.Rule,
				}
				if err := w.HandleRequest(ctx, msg); err != nil {
					respMsg.Success = false
					respMsg.Messages = []string{
						err.Error(),
					}
					log.Error().Err(err).Msg("error handling request")
				} else {
					respMsg.Success = true
				}

				if err := encoder.Encode(respMsg); err != nil {
					log.Error().Err(fmt.Errorf("failed to encode output: %w", err)).Msg("error handling request")
				}
			}()
		}
	}
}

func (w *PleaseWorker) HandleRequest(ctx context.Context, msg *Request) error {
	outImageOptValue, err := msg.ParseOption("--image_out")
	if err != nil {
		return fmt.Errorf("could not get option '--image_out': %w", err)
	}
	outImageFilePath := filepath.Join(msg.TempDir, outImageOptValue)
	outImageFile, err := os.Create(outImageFilePath)
	if err != nil {
		return fmt.Errorf("could not open '%s': %w", outImageFilePath, err)
	}
	defer outImageFile.Close()

	fqnTagsOptValue, err := msg.ParseOption("--fqn_tags_file")
	if err != nil {
		return fmt.Errorf("could not get option '--fqn_tags_file': %w", err)
	}
	fqnTagsFile := filepath.Join(msg.TempDir, fqnTagsOptValue)
	fqnTagsFileContents, err := os.ReadFile(fqnTagsFile)
	if err != nil {
		return fmt.Errorf("could not read '%s': %w", fqnTagsFile, err)
	}
	fqnTags := strings.FieldsFunc(string(fqnTagsFileContents), func(c rune) bool {
		return c == '\n'
	})

	// prepare dockerfile
	if err := os.MkdirAll(filepath.Join(msg.TempDir, "dockerfile"), 0755); err != nil {
		return fmt.Errorf("could not create 'dockerfile' dir: %w", err)
	}
	dockerfileOptValue, err := msg.ParseOption("--dockerfile")
	if err != nil {
		return fmt.Errorf("could not get option '--dockerfile': %w", err)
	}

	if err := os.Rename(
		filepath.Join(msg.TempDir, dockerfileOptValue),
		filepath.Join(msg.TempDir, "dockerfile/Dockerfile"),
	); err != nil {
		log.Printf("could not move dockerfile: %s", err)
	}

	cmd := exec.CommandContext(ctx,
		w.opts.BuildCtlBinary,
		[]string{
			"build",
			"--frontend=dockerfile.v0",
			"--local", fmt.Sprintf("context=%s", msg.TempDir),
			"--local", fmt.Sprintf("dockerfile=%s", filepath.Join(msg.TempDir, "dockerfile")),
			"--output", fmt.Sprintf("type=docker,\"name=%s\"", strings.Join(fqnTags, ",")),
		}...)
	cmd.Stdout = outImageFile
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Env = append(os.Environ(), []string{"BUILDKIT_HOST=" + w.opts.BuildKitAddress}...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not run '%s': %w\n%s", strings.Join(cmd.Args, " "), err, stderr.String())
	}

	return nil
}

// A Request is the message that's sent to a worker indicating that it should start a build.
type Request struct {
	// The label of the rule to build, i.e. //src/worker:worker
	Rule string `json:"rule"`
	// Labels applies to this rule.
	Labels []string `json:"labels"`
	// The temporary directory to build the target in.
	TempDir string `json:"temp_dir"`
	// List of source files to compile
	Sources []string `json:"srcs"`
	// Compiler options
	Options []string `json:"opts"`
	// True if this message relates to a test.
	Test bool `json:"test"`
}

// A Response is sent back from the worker on completion.
type Response struct {
	// The label of the rule to build, i.e. //src/worker:worker
	// Always corresponds to one that was sent out earlier in a request.
	Rule string `json:"rule"`
	// True if build succeeded
	Success bool `json:"success"`
	// Any messages reported. On failure these should indicate what's gone wrong.
	Messages []string `json:"messages"`
	// The contents of the BUILD file that should be assumed for this directory, if it's a parse request.
	BuildFile string `json:"build_file"`
	// If this is non-empty it replaces the existing test command.
	Command string `json:"command"`
}

func (r *Request) ParseOption(optionName string) (string, error) {
	for _, opt := range r.Options {
		optParts := strings.Split(opt, "=")
		if len(optParts) == 2 {
			optName := optParts[0]
			optValue := optParts[1]
			if optName == optionName {
				return strings.TrimPrefix(strings.TrimSuffix(optValue, "\""), "\""), nil
			}
		}
	}

	return "", fmt.Errorf("could not parse option '%s' from: %v", optionName, r.Options)
}
