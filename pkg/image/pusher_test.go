package image_test

import (
	"testing"

	"github.com/VJftw/please-buildkit/pkg/image"
	"github.com/stretchr/testify/assert"
)

func TestTranslateUserProvidedRepoTags(t *testing.T) {
	var tests = []struct {
		desc                   string
		inImageRepoTags        []string
		inUserProvidedRepoTags []string
		outRepoTags            []string
	}{
		{
			"img repository, img tag (no user input)",
			[]string{"img-registry.com/img-repository:img-tag"},
			[]string{},
			[]string{"img-registry.com/img-repository:img-tag"},
		},
		{
			"img repository, img tag (no user input) - multiple img tags",
			[]string{
				"img-registry.com/img-repository:img-tag",
				"img-registry.com/img-repository:img-tag-2",
			},
			[]string{},
			[]string{
				"img-registry.com/img-repository:img-tag",
				"img-registry.com/img-repository:img-tag-2",
			},
		},
		{
			"img repository, img tag (no user input) - multiple img domains",
			[]string{
				"img-registry.com/img-repository:img-tag",
				"img-registry-2.com/img-repository:img-tag",
			},
			[]string{},
			[]string{
				"img-registry.com/img-repository:img-tag",
				"img-registry-2.com/img-repository:img-tag",
			},
		},
		{
			"user-provided repository and tag (`usr-registry.com/usr-repository:usr-tag`)",
			[]string{"img-registry.com/img-repository:img-tag"},
			[]string{"usr-registry.com/usr-repository:usr-tag"},
			[]string{"usr-registry.com/usr-repository:usr-tag"},
		},
		{
			"user-provided repository and tag (`usr-registry.com/usr-repository:usr-tag`) - multiple img tags",
			[]string{
				"img-registry.com/img-repository:img-tag",
				"img-registry.com/img-repository:img-tag-2",
			},
			[]string{"usr-registry.com/usr-repository:usr-tag"},
			[]string{"usr-registry.com/usr-repository:usr-tag"},
		},
		{
			"user-provided repository and tag (`usr-registry.com/usr-repository:usr-tag`) - multiple img domains",
			[]string{
				"img-registry.com/img-repository:img-tag",
				"img-registry-2.com/img-repository:img-tag",
			},
			[]string{"usr-registry.com/usr-repository:usr-tag"},
			[]string{"usr-registry.com/usr-repository:usr-tag"},
		},
		{
			"user-provided repository and tag (`usr-registry.com/usr-repository:usr-tag usr-registry-2.com/usr-repository-2:usr-tag-2`) - multiple usr",
			[]string{"img-registry.com/img-repository:img-tag"},
			[]string{
				"usr-registry.com/usr-repository:usr-tag",
				"usr-registry-2.com/usr-repository-2:usr-tag-2",
			},
			[]string{
				"usr-registry.com/usr-repository:usr-tag",
				"usr-registry-2.com/usr-repository-2:usr-tag-2",
			},
		},
		{
			"user-provided repository, img tag (`usr-registry.com/usr-repository:`)",
			[]string{"img-registry.com/image:img-tag"},
			[]string{"usr-registry.com/usr-repository:"},
			[]string{"usr-registry.com/usr-repository:img-tag"},
		},
		{
			"user-provided repository, img tag (`usr-registry.com/usr-repository:`) - multiple img tags",
			[]string{
				"img-registry.com/img-repository:img-tag",
				"img-registry.com/img-repository:img-tag-2",
			},
			[]string{"usr-registry.com/usr-repository:"},
			[]string{
				"usr-registry.com/usr-repository:img-tag",
				"usr-registry.com/usr-repository:img-tag-2",
			},
		},
		{
			"user-provided repository, img tag (`usr-registry.com/usr-repository:`) - multiple img domains",
			[]string{
				"img-registry.com/img-repository:img-tag",
				"img-registry-2.com/img-repository:img-tag",
			},
			[]string{"usr-registry.com/usr-repository:"},
			[]string{"usr-registry.com/usr-repository:img-tag"},
		},
		{
			"user-provided repository, img tag (`usr-registry.com/usr-repository: usr-registry-2.com/usr-repository-2:`) - multiple usr repositories",
			[]string{"img-registry.com/image:img-tag"},
			[]string{
				"usr-registry.com/usr-repository:",
				"usr-registry-2.com/usr-repository-2:",
			},
			[]string{
				"usr-registry.com/usr-repository:img-tag",
				"usr-registry-2.com/usr-repository-2:img-tag",
			},
		},
		{
			"img repository, user-provided tag (`:usr-tag`)",
			[]string{"img-registry.com/image:img-tag"},
			[]string{":usr-tag"},
			[]string{"img-registry.com/image:usr-tag"},
		},
		{
			"img repository, user-provided tag (`:usr-tag`) - multiple img tags",
			[]string{
				"img-registry.com/img-repository:img-tag",
				"img-registry.com/img-repository:img-tag-2",
			},
			[]string{":usr-tag"},
			[]string{"img-registry.com/img-repository:usr-tag"},
		},
		{
			"img repository, user-provided tag (`:usr-tag`) - multiple img domains",
			[]string{
				"img-registry.com/img-repository:img-tag",
				"img-registry-2.com/img-repository:img-tag",
			},
			[]string{":usr-tag"},
			[]string{
				"img-registry.com/img-repository:usr-tag",
				"img-registry-2.com/img-repository:usr-tag",
			},
		},
		{
			"img repository, user-provided tag (`:usr-tag :usr-tag-2`) - multiple usr tags",
			[]string{"img-registry.com/image:img-tag"},
			[]string{
				":usr-tag",
				":usr-tag-2",
			},
			[]string{
				"img-registry.com/image:usr-tag",
				"img-registry.com/image:usr-tag-2",
			},
		},
		{
			"user-provided registry",
			[]string{"example.com/usr-image:usr-tag"},
			[]string{"usr-registry.com"},
			[]string{"usr-registry.com/usr-image:usr-tag"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			outRepoTags := image.TranslateUserProvidedRepoTags(
				tt.inImageRepoTags,
				tt.inUserProvidedRepoTags,
			)
			assert.ElementsMatch(t, tt.outRepoTags, outRepoTags)
		})
	}
}
