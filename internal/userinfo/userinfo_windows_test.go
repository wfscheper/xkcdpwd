package userinfo

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigDir(t *testing.T) {
	got := DefaultConfigDir("foo")
	if want, ok := os.LookupEnv("APPDATA"); ok {
		assert.Equal(t, want, got)
	} else {
		assert.Equal(t, "foo\\AppData\\Roaming", got)
	}
}

func TestDefaultConfigFile(t *testing.T) {
	if got, err := DefaultConfigFile("foo"); assert.NoError(t, err) {
		want := filepath.Join("C:", "/Users", os.Getenv("USER"), "AppData", "Roaming", "foo", "foo.conf")
		assert.Equal(t, want, got)
	}
}
