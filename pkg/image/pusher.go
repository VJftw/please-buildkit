package image

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

type PusherOpts struct {
	CraneTool   string
	FQNTagsPath string
}

type Pusher struct {
	opts *PusherOpts
}

func NewPusher(opts *PusherOpts) *Pusher {
	return &Pusher{
		opts: opts,
	}
}

func (p *Pusher) PushTar(ctx context.Context, tarPath string, repoTags []string) error {
	// TODO: return multiple errors when Go 1.20 is released.
	resErr := fmt.Errorf("")

	for _, repoTag := range repoTags {
		pushCmd := exec.CommandContext(ctx, p.opts.CraneTool, []string{
			"push",
			tarPath,
			repoTag,
		}...)

		pushCmd.Stderr = os.Stderr
		pushCmd.Stdout = os.Stdout

		if err := pushCmd.Run(); err != nil {
			log.Error().
				Err(err).
				Str("tarPath", tarPath).
				Str("repoTag", repoTag).
				Msg("could not push image")
			resErr = fmt.Errorf("%scould not push image: %s; ", resErr, err)
		}

		log.Info().
			Str("repoTag", repoTag).
			Msg("pushed image")
	}

	if resErr.Error() != "" {
		return resErr
	}

	return nil
}

// TranslateUserProvidedRepoTags translates user-provided repo tags into repo
// tags that are pushed:
//
// - Support pushing by user-provided repository and tag:
//
//	`push index.docker.io/foo/bar:my-tag`
//
// - Support pushing by user-provided repository, SBOM repository tags:
//
//	`push index.docker.io/foo/bar:`
//
// - Support pushing by SBOM repository, user-provided tag:
//
//	`push :my-tag`
//
// - Support pushing by user-provided registry, SBOM repository path and tags:
//
//	`push localhost:5000`
func TranslateUserProvidedRepoTags(imageRepoTags []string, userProvidedRepoTags []string) []string {
	if len(userProvidedRepoTags) < 1 {
		return imageRepoTags
	}

	imageRepos := []string{}
	imageTags := []string{}
	imagePaths := []string{}

	for _, irt := range imageRepoTags {
		repoTag := RepoTag(irt)

		imageRepos = append(imageRepos, repoTag.GetRepository())
		imageTags = append(imageTags, repoTag.GetTag())
		imagePaths = append(imagePaths, repoTag.GetPath())
	}

	translatedTagsSet := map[string]struct{}{}

	userProvidedRepoAndTagRegex := regexp.MustCompile(`^.+:?.+\/.+:.+`)
	userProvidedRegistryRegex := regexp.MustCompile(`^[^:/]+:?[^:]*$`)

	for _, uprt := range userProvidedRepoTags {
		switch {
		case uprt[0] == ':':
			// use the image repo, user-provided tag
			for _, repo := range imageRepos {
				translatedTagsSet[fmt.Sprintf("%s%s", repo, uprt)] = struct{}{}
			}
		case userProvidedRepoAndTagRegex.MatchString(uprt):
			// use the user-provided repo and tag
			translatedTagsSet[uprt] = struct{}{}
		case uprt[len(uprt)-1] == ':':
			// use the user-provided repo, image tag
			for _, tag := range imageTags {
				translatedTagsSet[fmt.Sprintf("%s%s", uprt, tag)] = struct{}{}
			}
		case userProvidedRegistryRegex.MatchString(uprt):
			// use the user-provided registry, image path and tags.
			for _, imgPath := range imagePaths {
				for _, imgTag := range imageTags {
					translatedTagsSet[fmt.Sprintf("%s/%s:%s", uprt, imgPath, imgTag)] = struct{}{}
				}
			}
		}
	}

	translatedTags := []string{}
	for translatedTag := range translatedTagsSet {
		translatedTags = append(translatedTags, translatedTag)
	}

	sort.Strings(translatedTags)

	return translatedTags
}

func LoadImageRepoTags(fqnTagsPath string) ([]string, error) {
	fqnTagsFileBytes, err := os.ReadFile(fqnTagsPath)
	if err != nil {
		return nil, fmt.Errorf("could not read '%s': %w", fqnTagsPath, err)
	}

	fqnTags := strings.FieldsFunc(
		string(fqnTagsFileBytes),
		func(c rune) bool { return c == '\n' },
	)

	return fqnTags, nil
}
