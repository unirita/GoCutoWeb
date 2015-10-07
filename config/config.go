// Package config implements access for config file.
package config

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
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

	Server = c.Server
	Log = c.Log

	return nil
}
