package currencyconverter

import (
	"testing"
	"time"
)

var (
	testFileCacheEnabled        = false
	testFileCacheFilenameSuffix = "history"
	testFileCacheTimeout        = 30 * time.Minute
)

func TestFileCache(t *testing.T) {
	c := NewFileCache(FileCache{})

	c.SetEnabled(testFileCacheEnabled)
	if c.Enabled != testFileCacheEnabled {
		t.Error("Expected result is", testFileCacheEnabled, "but got", c.Enabled)
	}

	c.SetFilenameSuffix(testFileCacheFilenameSuffix)
	if c.FilenameSuffix != testFileCacheFilenameSuffix {
		t.Error("Expected result is", testFileCacheFilenameSuffix, "but got", c.FilenameSuffix)
	}

	c.SetTimeout(testFileCacheTimeout)
	if c.Timeout != testFileCacheTimeout {
		t.Error("Expected result is", testFileCacheTimeout, "but got", c.Timeout)
	}
}
