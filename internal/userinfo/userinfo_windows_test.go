package userinfo

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigDir(t *testing.T) {
	got := DefaultConfigDir("foo")
	assert.Equal(t, "foo\\Library\\Application Support", got)
}

func TestDefaultConfigFile(t *testing.T) {
	if got, err := DefaultConfigFile("foo"); assert.NoError(t, err) {
		want := filepath.Join("/Users", os.Getenv("USER"), "Library/Application Support", "foo", "foo.conf")
		assert.Equal(t, want, got)
	}
}
