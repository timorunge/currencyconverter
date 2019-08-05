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
// FileCacheFilename is the default name for the file cache file - without file
// extention.
// FileCacheTimeout is defining the timeout for the cache file.
const (
	FileCacheEnabled  = true
	FileCacheFilename = "currencyconverter"
	FileCacheTimeout  = 60 * time.Minute
)

// FileCacheDirectory is defining the path where to store the cache file.
var (
	FileCacheDirectory = os.TempDir()
)

// FileCache is the struct for the file cache.
type FileCache struct {
	Directory     string
	Enabled       bool
	ExchangeRates ExchangeRates
	Filename      string
	FullFilePath  string
	Timeout       time.Duration
}

// NewFileCache returning a new file cache struct.
func NewFileCache() *FileCache {
	c := &FileCache{
		Directory: FileCacheDirectory,
		Enabled:   FileCacheEnabled,
		Filename:  FileCacheFilename,
		Timeout:   FileCacheTimeout,
	}
	return c.setFullFilePath()
}

// IsValidDirectory is checking if the path is valid.
func IsValidDirectory(directory string) error {
	c := &FileCache{Directory: directory}
	return c.isValidDirectory()
}

// IsValidDirectory is checking if the path is valid.
func (c *FileCache) IsValidDirectory() error {
	return c.isValidDirectory()
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

// SetDirectory is setting the directory to store the file cache.
func (c *FileCache) SetDirectory(directory string) *FileCache {
	c.Directory = directory
	return c.setFullFilePath()
}

// SetEnabled is enabling or disabling the file cache.
func (c *FileCache) SetEnabled(enabled bool) *FileCache {
	c.Enabled = enabled
	return c
}

// SetFilename is setting the name of the file cache file. File extension will
// automatically be added.
func (c *FileCache) SetFilename(filename string) *FileCache {
	c.Filename = filename
	return c.setFullFilePath()
}

// SetTimeout is setting the the file cache timeout.
func (c *FileCache) SetTimeout(timeout time.Duration) *FileCache {
	c.Timeout = timeout
	return c
}

// isValidDirectory is checking if the path is valid.
func (c *FileCache) isValidDirectory() error {
	directory, err := os.Stat(c.Directory)
	if err != nil {
		return fmt.Errorf("Can not find directory \"%s\"", c.Directory)
	}
	if !directory.IsDir() {
		return fmt.Errorf("\"%s\" is an existing file, not a directory", c.Directory)
	}
	return nil
}

// getExchangeRates is getting the exchange rates from the cache.
func (c *FileCache) getExchangeRates() (ExchangeRates, error) {
	c.setFullFilePath()
	file, err := os.Stat(c.FullFilePath)
	if err != nil {
		return c.ExchangeRates, err
	}
	if time.Duration(c.Timeout) > time.Now().Sub(file.ModTime()) {
		fileContent, err := ioutil.ReadFile(c.FullFilePath)
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
	c.setFullFilePath()
	file, err := os.Create(c.FullFilePath)
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

// setFullFilePath is setting the complete path for the cache file.
func (c *FileCache) setFullFilePath() *FileCache {
	c.FullFilePath = fmt.Sprintf("%s/%s.json", c.Directory, c.Filename)
	return c
}
