// JSON file-based caching for exchange rates.

package converter

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	errFileCacheDisabled = errors.New("filecache: caching is disabled")
	errFileCacheExpired  = errors.New("filecache: cache expired")
	errFileCacheNoRate   = errors.New("filecache: no exchange rate provided")
)

// FileCacheConfig holds configuration for a FileCache.
type FileCacheConfig struct {
	Directory string
	Enabled   bool
	Filename  string
	MaxAge    time.Duration // Zero means no expiry.
}

// FileCache provides JSON file-based caching for exchange rates.
type FileCache struct {
	config FileCacheConfig
}

// NewFileCache returns a new FileCache.
func NewFileCache(cfg FileCacheConfig) *FileCache {
	if cfg.Directory == "" {
		cfg.Directory = os.TempDir()
	} else {
		cfg.Directory = filepath.Clean(cfg.Directory)
	}
	cfg.Filename = cmp.Or(cfg.Filename, "currencyconverter")
	cfg.Filename = filepath.Base(cfg.Filename)
	return &FileCache{config: cfg}
}

// Get returns cached DailyRates, or an error if the cache is
// disabled, missing, or expired.
func (c *FileCache) Get() (DailyRates, error) {
	if !c.config.Enabled {
		return DailyRates{}, errFileCacheDisabled
	}
	path := c.filePath()
	f, err := os.Open(path)
	if err != nil {
		return DailyRates{}, err
	}
	defer func() { _ = f.Close() }()

	if c.config.MaxAge > 0 {
		info, statErr := f.Stat()
		if statErr != nil {
			return DailyRates{}, statErr
		}
		if time.Since(info.ModTime()) > c.config.MaxAge {
			return DailyRates{}, errFileCacheExpired
		}
	}

	var d DailyRates
	if err := json.NewDecoder(f).Decode(&d); err != nil {
		return DailyRates{}, fmt.Errorf("filecache: %w", err)
	}
	if len(d.Rates) == 0 {
		return DailyRates{}, errFileCacheNoRate
	}
	return d, nil
}

// Write persists DailyRates to the cache file atomically.
func (c *FileCache) Write(date DailyRates) error {
	tmp, err := os.CreateTemp(c.config.Directory, c.config.Filename+"*.tmp")
	if err != nil {
		return fmt.Errorf("filecache: %w", err)
	}
	tmpPath := tmp.Name()

	if err := json.NewEncoder(tmp).Encode(date); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("filecache: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("filecache: %w", err)
	}
	if err := os.Rename(tmpPath, c.filePath()); err != nil {
		return fmt.Errorf("filecache: %w", err)
	}
	return nil
}

// IsCacheMiss reports whether the error represents an expected cache
// miss (disabled, expired, or file not found) as opposed to corruption.
func IsCacheMiss(err error) bool {
	if errors.Is(err, errFileCacheDisabled) || errors.Is(err, errFileCacheExpired) || errors.Is(err, errFileCacheNoRate) {
		return true
	}
	return os.IsNotExist(err)
}

func (c *FileCache) filePath() string {
	return filepath.Join(c.config.Directory, c.config.Filename+".json")
}
