package currencyconverter

import (
	"fmt"
	"testing"
	"time"
)

func TestFileCacheDirectory(t *testing.T) {
	testFileCacheDirectory := "/tmp"

	c := NewFileCache()
	c.SetDirectory(testFileCacheDirectory)
	if c.Directory != testFileCacheDirectory {
		t.Error("Expected result is", testFileCacheDirectory, "but got", c.Directory)
	}
}

func TestFileCacheEnabled(t *testing.T) {
	testFileCacheEnabled := false

	c := NewFileCache()
	c.SetEnabled(testFileCacheEnabled)
	if c.Enabled != testFileCacheEnabled {
		t.Error("Expected result is", testFileCacheEnabled, "but got", c.Enabled)
	}
}

func TestFileCacheFilename(t *testing.T) {
	testFileCacheFilename := "today"

	c := NewFileCache()
	c.SetFilename(testFileCacheFilename)
	if c.Filename != testFileCacheFilename {
		t.Error("Expected result is", testFileCacheFilename, "but got", c.Filename)
	}
}

func TestFileCacheTimeout(t *testing.T) {
	testFileCacheTimeout := 30 * time.Minute

	c := NewFileCache()
	c.SetTimeout(testFileCacheTimeout)
	if c.Timeout != testFileCacheTimeout {
		t.Error("Expected result is", testFileCacheTimeout, "but got", c.Timeout)
	}
}

func TestFileCacheFullFilePath(t *testing.T) {
	testFileCacheDirectory := "/tmp"
	testFileCacheFilename := "today"
	testFileCacheFullFilePath := fmt.Sprintf("%s/%s.json", testFileCacheDirectory, testFileCacheFilename)

	c := NewFileCache()
	c.SetDirectory(testFileCacheDirectory)
	c.SetFilename(testFileCacheFilename)
	if c.FullFilePath != testFileCacheFullFilePath {
		t.Error("Expected result is", testFileCacheFullFilePath, "but got", c.FullFilePath)
	}
}
