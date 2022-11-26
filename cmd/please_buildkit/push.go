package main

import (
	"fmt"

	"github.com/VJftw/please-buildkit/pkg/image"
	"github.com/urfave/cli/v2"
)

func PushCommand() *cli.Command {
	return &cli.Command{
		Name:  "push",
		Usage: "Pushes the given image tar",
		Description: `
This command pushes the given 'image_tar_path' to the repositories using the
given 'crane_tool' and 'fqn_tags_path' it was created with. You may provide arguments to
augment the repositories and tags that get pushed. This supports the following
augmentations as arguments:

 - ` + "`index.docker.io/foo/bar:other-tag`\t" + `push to user-provided repository and tag.
 - ` + "`index.docker.io/foo/bar:`\t" + `push to user-provided repository with image-defined tags. Note the trailing ':'.
 - ` + "`:other-tag`\t" + `push to image-defined repository with user-provided tags. Note the leading ':'.
 - ` + "`localhost:5000`\t" + `push to user-provided registry with image-defined repository and tags.

This maintains consistency with the 'replace' command so that it is easy to use
the same arguments with both commands.

For example, given an image tar tagged with 'index.docker.io/my-repo:my-tag':

  1. ` + "`" + `$ push "index.docker.io/other-repo:other-tag"` + "`" + `:
     Will push the image to only "index.docker.io/other-repo:other-tag", regardless
     of other image tags.
  2. ` + "`" + `$ push "index.docker.io/foo/bar:"` + "`" + `:
  	 Will push the image to "index.docker.io/foo/bar:my-tag".
  3. ` + "`" + `$ push ":other-tag"` + "`" + `:
  	 Will push the image to "index.docker.io/my-repo:other-tag".
  4. ` + "`" + `$ push "localhost:5000"` + "`" + `:
     Will push the image to "localhost:5000/my-repo:my-tag".
`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "crane_tool",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "img_tar_path",
				Required: true,
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

			imageRepoTagsToPush := image.TranslateUserProvidedRepoTags(imageTags, userProviderRepoTags)

			imagePusher := image.NewPusher(&image.PusherOpts{
				CraneTool: cCtx.String("crane_tool"),
			})

			return imagePusher.PushTar(cCtx.Context, cCtx.String("img_tar_path"), imageRepoTagsToPush)
		},
	}
}
