package main

import (
	"fmt"
	"os"

	"github.com/VJftw/please-buildkit/pkg/image"
	"github.com/urfave/cli/v2"
)

func ReplaceCommand() *cli.Command {
	return &cli.Command{
		Name:  "replace",
		Usage: "Replaces references to an image in the given file path",
		Description: `
This command replaces references to an image defined by the 'aliases' in the
given file path with a reference from the 'fqn_tags_path'. This supports the following
augmentations to the replaced image as an argument. Note only 1 augmentation is
allowed per invocation:

 - ` + "`index.docker.io/foo/bar:other-tag`\t" + `replace to user-provided repository and tag.
 - ` + "`index.docker.io/foo/bar:`\t" + `replace to user-provided repository with image-defined tag. Note the trailing ':'.
 - ` + "`:other-tag`\t" + `replace to image-defined repository with user-provided tags. Note the leading ':'.
 - ` + "`localhost:5000`\t" + `replace to user-provided registry with image-defined repository and tags.

These augmentations are consistent with the 'push' command so that it is easy to
use the same arguments with both commands.

During replacement, this command prioritises tags in the following order:
  1. ":srcsha256-*".
  2. Any non-latest tag.
  3. ":latest".
`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file_path",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name: "aliases",
			},
			&cli.StringFlag{
				Name:     "fqn_tags_path",
				Required: true,
			},
		},
		Action: func(cCtx *cli.Context) error {
			userProviderRepoTags := cCtx.Args().Slice()

			imageTags, err := image.LoadImageRepoTags(cCtx.String("fqn_tags_path"))
			if err != nil {
				return fmt.Errorf("could not load image tags: %w", err)
			}

			imageRepoTagsToReplace := image.TranslateUserProvidedRepoTags(imageTags, userProviderRepoTags)
			granularestRepoTagToReplace := image.DetermineMostGranularRepoTags(imageRepoTagsToReplace)

			targetPath := cCtx.String("file_path")
			contents, err := os.ReadFile(targetPath)
			if err != nil {
				return fmt.Errorf("could not read '%s': %w", cCtx.String("fqn_tags_path"), err)
			}

			for _, alias := range cCtx.StringSlice("aliases") {
				contents = image.ReplaceImageReferences(contents, alias, granularestRepoTagToReplace)
			}

			stat, err := os.Stat(targetPath)
			if err != nil {
				return fmt.Errorf("could not stat '%s': %w", targetPath, err)
			}

			if err := os.WriteFile(targetPath, contents, stat.Mode()); err != nil {
				return fmt.Errorf("could not write '%s': %w", targetPath, err)
			}

			return nil
		},
	}
}
