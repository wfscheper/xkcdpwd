package userinfo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigDir(t *testing.T) {
	actual := DefaultConfigDir("foo")
	if "foo/Library/Application Support" != actual {
		t.Errorf("wrong config dir")
	}
}

func TestDefaultConfigFile(t *testing.T) {
	actual, err := DefaultConfigFile("foo")
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join("/Users", os.Getenv("USER"), "Library/Application Support", "foo", "foo.conf")
	if actual != expected {
		t.Errorf("wrong config file")
	}
}
