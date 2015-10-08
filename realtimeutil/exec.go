package realtimeutil

import (
	"bytes"
	"os/exec"
	"path/filepath"

	"github.com/unirita/gocutoweb/config"
	"github.com/unirita/gocutoweb/log"
)

type Command struct {
	cmd    *exec.Cmd
	Result string
}

func NewCommand(dynamicJobnetName, url string) *Command {
	realtime := filepath.Join(config.Server.MasterDir, "realtime")

	c := new(Command)
	c.cmd = exec.Command(realtime, "-n", dynamicJobnetName, url)
	return c
}

// Run runs the realtime-utility command.
func (c *Command) Run() error {
	var stdout bytes.Buffer
	c.cmd.Stdout = &stdout
	c.cmd.Stderr = &stdout

	if err := c.cmd.Start(); err != nil {
		log.Error("Execute realtime-utility failed - %v", err.Error())
		return err
	}
	c.cmd.Wait()
	c.Result = stdout.String()
	return nil
}
