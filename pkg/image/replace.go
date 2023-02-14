package image

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrNoReplacementsMade = errors.New("no replacements made")
)

func DetermineMostGranularRepoTags(repoTags []string) string {
	var currentGranularestRepoTag string
	for _, rT := range repoTags {
		repoTag := RepoTag(rT)
		tag := repoTag.GetTag()
		if strings.HasPrefix(tag, "srcsha256-") {
			return rT
		}

		if tag != "latest" && currentGranularestRepoTag == "" {
			currentGranularestRepoTag = rT
		}
	}

	if currentGranularestRepoTag != "" {
		return currentGranularestRepoTag
	}

	return repoTags[0]
}

func ReplaceImageReferences(contents []byte, oldRef string, newRef string) ([]byte, error) {

	oldRepoTag := RepoTag(oldRef)
	var oldRefRegex *regexp.Regexp
	if oldRepoTag.GetTag() == "" {
		// not fq
		oldRefOnlyRegex := regexp.QuoteMeta(oldRef)
		regex := fmt.Sprintf(`%s[a-z0-9:@\.\-\_]*`, oldRefOnlyRegex)
		oldRefRegex = regexp.MustCompile(regex)

		if len(oldRefRegex.FindAll(contents, -1)) < 1 {
			return contents, fmt.Errorf("could not replace not fully-qualified image '%s': %w", regex, ErrNoReplacementsMade)
		}

		return oldRefRegex.ReplaceAll(contents, []byte(newRef)), nil
	}

	if !bytes.Contains(contents, []byte(oldRef)) {
		return contents, fmt.Errorf("could not replace fully-qualified image '%s': %w", oldRef, ErrNoReplacementsMade)

	}

	return bytes.ReplaceAll(contents, []byte(oldRef), []byte(newRef)), nil
}

func ReplaceImageReferencesForAliases(contents []byte, newRef string, aliases ...string) ([]byte, error) {
	replaced := false
	var errs error
	for _, alias := range aliases {
		newContents, err := ReplaceImageReferences(contents, alias, newRef)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("replacing alias '%s' encountered error: %w", alias, err))
			continue
		}
		replaced = true
		contents = newContents
	}

	if !replaced {
		return contents, errs
	}

	return contents, nil
}
