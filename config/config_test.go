package config

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	config := `
[server]
listen_port = 1234
master_path = "/usr/local/gocuto/bin/master"

[log]
output_dir     = "/var/log/gocuto"
output_level   = "info"
max_size_kb    = 10240
max_generation = 2
`
	err := loadReader(strings.NewReader(config))
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if Server.ListenPort != 1234 {
		t.Errorf("server.listen_port => %d, want %d", Server.ListenPort, 1234)
	}
	if Server.MasterPath != "/usr/local/gocuto/bin/master" {
		t.Errorf("server.master_path => %s, want %s", Server.MasterPath, "/usr/local/gocuto/bin/master")
	}
	if Log.OutputDir != "/var/log/gocuto" {
		t.Errorf("log.output_dir => %s, want %s", Log.OutputDir, "/var/log/gocuto")
	}
	if Log.OutputLevel != "info" {
		t.Errorf("log.output_level => %s, want %s", Log.OutputLevel, "info")
	}
	if Log.MaxSizeKB != 10240 {
		t.Errorf("log.max_size_kb => %d, want %d", Log.MaxSizeKB, 10240)
	}
	if Log.MaxGeneration != 2 {
		t.Errorf("log.max_generation => %d, want %d", Log.MaxGeneration, 2)
	}
}
