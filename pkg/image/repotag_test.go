package image_test

import (
	"testing"

	"github.com/VJftw/please-buildkit/pkg/image"
	"github.com/stretchr/testify/assert"
)

func TestRepoTagGetRepository(t *testing.T) {
	var tests = []struct {
		image         image.RepoTag
		outRepository string
	}{
		{image.RepoTag("example.com/foo:bar"), "example.com/foo"},
		{image.RepoTag("localhost:5000/foo:bar"), "localhost:5000/foo"},
	}

	for _, tt := range tests {
		t.Run(string(tt.image), func(t *testing.T) {
			repo := tt.image.GetRepository()
			assert.Equal(t, tt.outRepository, repo)
		})
	}
}

func TestRepoTagGetTag(t *testing.T) {
	var tests = []struct {
		image  image.RepoTag
		outTag string
	}{
		{image.RepoTag("example.com/foo:bar"), "bar"},
		{image.RepoTag("localhost:5000/foo:bar"), "bar"},
	}

	for _, tt := range tests {
		t.Run(string(tt.image), func(t *testing.T) {
			repo := tt.image.GetTag()
			assert.Equal(t, tt.outTag, repo)
		})
	}
}

func TestRepoTagGetRegistry(t *testing.T) {
	var tests = []struct {
		image       image.RepoTag
		outRegistry string
	}{
		{image.RepoTag("example.com/foo:bar"), "example.com"},
		{image.RepoTag("localhost:5000/foo:bar"), "localhost:5000"},
	}

	for _, tt := range tests {
		t.Run(string(tt.image), func(t *testing.T) {
			repo := tt.image.GetRegistry()
			assert.Equal(t, tt.outRegistry, repo)
		})
	}
}

func TestRepoTagGetPath(t *testing.T) {
	var tests = []struct {
		image   image.RepoTag
		outPath string
	}{
		{image.RepoTag("example.com/foo:bar"), "foo"},
		{image.RepoTag("localhost:5000/foo:bar"), "foo"},
	}

	for _, tt := range tests {
		t.Run(string(tt.image), func(t *testing.T) {
			repo := tt.image.GetPath()
			assert.Equal(t, tt.outPath, repo)
		})
	}
}
