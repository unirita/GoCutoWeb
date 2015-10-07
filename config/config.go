// Package config implements access for config file.
package config

import (
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/unirita/gocutoweb/pathutil"
)

// Config holds all config parameter.
type Config struct {
	Server serverSection
	Log    logSection
}

type serverSection struct {
	ListenPort int    `toml:"listen_port"`
	MasterPath string `toml:"master_path"`
}

type logSection struct {
	OutputDir     string `toml:"output_dir"`
	OutputLevel   string `toml:"output_level"`
	MaxSizeKB     int    `toml:"max_size_kb"`
	MaxGeneration int    `toml:"max_generation"`
}

const tag_CUTOROOT = "<CUTOROOT>"

var Server serverSection
var Log logSection

// Load loads config file from path, and returns config object as singleton.
//
// Caution: this function is not thread safe.
func Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	return loadReader(f)
}

func loadReader(r io.Reader) error {
	c := new(Config)
	if _, err := toml.DecodeReader(r, c); err != nil {
		return err
	}

	replaceCutoroot(c)

	Server = c.Server
	Log = c.Log

	return nil
}

func replaceCutoroot(c *Config) {
	c.Server.MasterPath = strings.Replace(c.Server.MasterPath, tag_CUTOROOT, pathutil.GetRootPath(), -1)
	c.Log.OutputDir = strings.Replace(c.Log.OutputDir, tag_CUTOROOT, pathutil.GetRootPath(), -1)
}
