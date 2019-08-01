package currencyconverter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// FileCacheEnabled is enabling / disabling the file cache.
// FileCacheFilenameSuffix is the default suffix for the file cache name.
// FileCacheTimeout is defining the timeout for the cache file.
const (
	FileCacheEnabled        = true
	FileCacheFilenameSuffix = "latest"
	FileCacheTimeout        = 60 * time.Minute
)

// FileCache is the struct for the file cache.
type FileCache struct {
	Enabled        bool
	Filename       string
	FilenameSuffix string
	ExchangeRates  ExchangeRates
	Timeout        time.Duration
}

// NewFileCache returning a new file cache struct.
func NewFileCache(cache FileCache) *FileCache {
	c := &FileCache{
		Enabled:        cache.Enabled,
		FilenameSuffix: cache.FilenameSuffix,
		ExchangeRates:  cache.ExchangeRates,
		Timeout:        cache.Timeout,
	}
	c.setFilename()
	return c
}

// Get is returning the exchange rates from the file cache.
func (c *FileCache) Get() (ExchangeRates, error) {
	if c.Enabled {
		return c.getExchangeRates()
	}
	return c.ExchangeRates, errors.New("Caching is disabled")
}

// Write is writing exchange rates to the file cache.
func (c *FileCache) Write(exchangeRates ExchangeRates) error {
	return c.writeToFileCache(exchangeRates)
}

// SetEnabled is enabling or disabling the file cache.
func (c *FileCache) SetEnabled(enabled bool) {
	c.Enabled = enabled
}

// SetFilenameSuffix is setting the suffix for the cache file name.
func (c *FileCache) SetFilenameSuffix(suffix string) {
	c.setFilenameSuffix(suffix)
}

// SetTimeout is setting the the file cache timeout.
func (c *FileCache) SetTimeout(timeout time.Duration) {
	c.Timeout = timeout
}

// getExchangeRates is getting the exchange rates from the cache.
func (c *FileCache) getExchangeRates() (ExchangeRates, error) {
	c.setFilename()
	cacheFile, err := os.Stat(c.Filename)
	if err != nil {
		return c.ExchangeRates, err
	}
	if time.Duration(c.Timeout) > time.Now().Sub(cacheFile.ModTime()) {
		fileContent, err := ioutil.ReadFile(c.Filename)
		if err != nil {
			return c.ExchangeRates, err
		}
		if err := json.Unmarshal(fileContent, &c.ExchangeRates); err != nil {
			return c.ExchangeRates, err
		}
		if len(c.ExchangeRates.Dates) == 0 {
			return c.ExchangeRates, errors.New("Not getting a single exchange rate from the file cache")
		}
		return c.ExchangeRates, nil
	}
	return c.ExchangeRates, errors.New("File cache is outdated")
}

// writeToFileCache is writing exchange rates to the file cache.
func (c *FileCache) writeToFileCache(exchangeRates ExchangeRates) error {
	c.setFilename()
	file, err := os.Create(c.Filename)
	defer file.Close()
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(exchangeRates)
	if err != nil {
		return err
	}
	if _, err := file.Write(bytes); err != nil {
		return err
	}
	return nil
}

// setFilename is setting the cache file filename.
func (c *FileCache) setFilename() {
	if c.FilenameSuffix == "" {
		c.FilenameSuffix = FileCacheFilenameSuffix
	}
	c.Filename = fmt.Sprintf("%s/currencyconverter-%s.json", os.TempDir(), c.FilenameSuffix)
}

// setFilenameSuffix is setting the suffix for the cache file filename.
func (c *FileCache) setFilenameSuffix(suffix string) {
	c.FilenameSuffix = suffix
	c.setFilename()
}
