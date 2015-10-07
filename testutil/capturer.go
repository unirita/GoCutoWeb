package testutil

import (
	"bytes"
	"io"
	"os"
)

// Capturer is object to capture stdout/stderr.
type Capturer struct {
	isStderr bool
	original *os.File
	bufCh    chan string
	out      *os.File
	in       *os.File
}

// NewStderrCapturer creates capturer for stdout.
func NewStdoutCapturer() *Capturer {
	c := new(Capturer)
	c.isStderr = false
	return c
}

// NewStderrCapturer creates capturer for stderr.
func NewStderrCapturer() *Capturer {
	c := new(Capturer)
	c.isStderr = true
	return c
}

// Start starts capturing.
func (c *Capturer) Start() {
	if c.isStderr {
		c.original = os.Stderr
	} else {
		c.original = os.Stdout
	}
	var err error
	c.in, c.out, err = os.Pipe()
	if err != nil {
		panic(err)
	}

	if c.isStderr {
		os.Stderr = c.out
	} else {
		os.Stdout = c.out
	}
	c.bufCh = make(chan string)
	go func() {
		var b bytes.Buffer
		io.Copy(&b, c.in)
		c.bufCh <- b.String()
	}()
}

// Stop stops capturing.
func (c *Capturer) Stop() string {
	c.out.Close()
	defer c.in.Close()
	if c.isStderr {
		os.Stderr = c.original
	} else {
		os.Stdout = c.original
	}
	return <-c.bufCh
}
