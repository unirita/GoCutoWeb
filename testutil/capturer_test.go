package testutil

import (
	"fmt"
	"os"
	"testing"
)

func TestCapturer_stdout(t *testing.T) {
	c := NewStdoutCapturer()
	c.Start()
	fmt.Println("test")
	output := c.Stop()
	if output != "test\n" {
		t.Error("Capture result is not expected value.")
		t.Logf("Expected: %s", "test\n")
		t.Logf("Actual: %s", output)
	}
}

func TestCapturer_stderr(t *testing.T) {
	c := NewStderrCapturer()
	c.Start()
	fmt.Fprintln(os.Stderr, "test")
	output := c.Stop()
	if output != "test\n" {
		t.Error("Capture result is not expected value.")
		t.Logf("Expected: %s", "test\n")
		t.Logf("Actual: %s", output)
	}
}
