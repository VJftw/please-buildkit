package image_test

import (
	"testing"

	"github.com/VJftw/please-buildkit/pkg/image"
	"github.com/stretchr/testify/assert"
)

func TestDetermineMostGranularRepoTag(t *testing.T) {
	var tests = []struct {
		desc               string
		inImageRepoTags    []string
		outGranularRepoTag string
	}{
		{
			"prioritise :srcsha256-* over :latest",
			[]string{"registry.com/repo:latest", "registry.com/repo:srcsha256-12345"},
			"registry.com/repo:srcsha256-12345",
		},
		{
			"prioritise :srcsha256-* over :other-tag",
			[]string{"registry.com/repo:other-tag", "registry.com/repo:srcsha256-12345"},
			"registry.com/repo:srcsha256-12345",
		},
		{
			"prioritise :other-tag over :latest",
			[]string{"registry.com/repo:latest", "registry.com/repo:other-tag"},
			"registry.com/repo:other-tag",
		},
		{
			"fallback to :latest",
			[]string{"registry.com/repo:latest"},
			"registry.com/repo:latest",
		},
		{
			"prioritise :srcsha256-* over :latest - multiple repositories",
			[]string{
				"registry.com/repo:latest",
				"registry.com/repo:srcsha256-12345",
				"registry.com/repo-2:latest",
				"registry.com/repo-2:srcsha256-12345",
			},
			"registry.com/repo:srcsha256-12345",
		},
		{
			"prioritise :srcsha256-* over :other-tag - multiple repositories",
			[]string{
				"registry.com/repo:other-tag",
				"registry.com/repo:srcsha256-12345",
				"registry.com/repo-2:other-tag",
				"registry.com/repo-2:srcsha256-12345",
			},
			"registry.com/repo:srcsha256-12345",
		},
		{
			"prioritise :other-tag over :latest - multiple repositories",
			[]string{
				"registry.com/repo:other-tag",
				"registry.com/repo:latest",
				"registry.com/repo-2:other-tag",
				"registry.com/repo-2:latest",
			},
			"registry.com/repo:other-tag",
		},
		{
			"prioritise :srcsha256-* over :latest - multiple repositories and tags",
			[]string{
				"registry.com/repo:latest",
				"registry.com/repo-2:latest",
				"registry.com/repo-2:srcsha256-12345",
			},
			"registry.com/repo-2:srcsha256-12345",
		},
		{
			"prioritise :srcsha256-* over :other-tag - multiple repositories and tags",
			[]string{
				"registry.com/repo:other-tag",
				"registry.com/repo-2:other-tag",
				"registry.com/repo-2:srcsha256-12345",
			},
			"registry.com/repo-2:srcsha256-12345",
		},
		{
			"prioritise :other-tag over :latest - multiple repositories and tags",
			[]string{
				"registry.com/repo:latest",
				"registry.com/repo-2:other-tag",
				"registry.com/repo-2:latest",
			},
			"registry.com/repo-2:other-tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			outGranularRepoTag := image.DetermineMostGranularRepoTags(
				tt.inImageRepoTags,
			)
			assert.Equal(t, tt.outGranularRepoTag, outGranularRepoTag)
		})
	}
}

func TestReplaceImageReferences(t *testing.T) {
	var tests = []struct {
		desc        string
		inContents  string
		inOldRef    string
		inNewRef    string
		outContents string
	}{
		{
			"given repo tag, no quotes in contents",
			`image: registry.com/foo/bar:latest`,
			`registry.com/foo/bar:latest`,
			`registry.com/foo/bar:srcsha-12345`,
			`image: registry.com/foo/bar:srcsha-12345`,
		},
		{
			"given repo tag, quotes in contents",
			`image: "registry.com/foo/bar:latest"`,
			`registry.com/foo/bar:latest`,
			`registry.com/foo/bar:srcsha-12345`,
			`image: "registry.com/foo/bar:srcsha-12345"`,
		},
		{
			"given repo, no quotes in contents",
			`image: registry.com/foo/bar:latest`,
			`registry.com/foo/bar`,
			`registry.com/foo/bar:srcsha-12345`,
			`image: registry.com/foo/bar:srcsha-12345`,
		},
		{
			"given repo, quotes in contents",
			`image: "registry.com/foo/bar:latest"`,
			`registry.com/foo/bar`,
			`registry.com/foo/bar:srcsha-12345`,
			`image: "registry.com/foo/bar:srcsha-12345"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			outContents := image.ReplaceImageReferences(
				[]byte(tt.inContents),
				tt.inOldRef,
				tt.inNewRef,
			)

			assert.Equal(t, tt.outContents, string(outContents))
		})
	}

}
