package image

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
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

func ReplaceImageReferences(contents []byte, oldRef string, newRef string) []byte {

	oldRepoTag := RepoTag(oldRef)
	var oldRefRegex *regexp.Regexp
	if oldRepoTag.GetTag() == "" {
		// not fq
		oldRefOnlyRegex := regexp.QuoteMeta(oldRef)
		oldRefRegex = regexp.MustCompile(fmt.Sprintf(`%s[a-z1-9:@\.\-\_]+`, oldRefOnlyRegex))

		return oldRefRegex.ReplaceAll(contents, []byte(newRef))
	}

	return bytes.ReplaceAll(contents, []byte(oldRef), []byte(newRef))
}
