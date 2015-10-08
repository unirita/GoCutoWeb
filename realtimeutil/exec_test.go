package realtimeutil

import (
	"testing"
)

var testJobnet string = "testjobnet"
var testURL string = "http://127.0.0.1"

func TestNewCommand(t *testing.T) {
	c := NewCommand(testJobnet, testURL)

	if c.cmd.Args[0] != "realtime" {
		t.Errorf("c.cmd.Args[0] -> %v, want %v", c.cmd.Args[0], "realtime")
	}
	if c.cmd.Args[1] != "-n" {
		t.Errorf("c.cmd.Args[0] -> %v, want %v", c.cmd.Args[1], "-n")
	}
	if c.cmd.Args[2] != testJobnet {
		t.Errorf("c.cmd.Args[0] -> %v, want %v", c.cmd.Args[2], testJobnet)
	}
	if c.cmd.Args[3] != testURL {
		t.Errorf("c.cmd.Args[0] -> %v, want %v", c.cmd.Args[3], testURL)
	}
}
