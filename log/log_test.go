package log

import (
	"sync"
	"testing"

	"github.com/cihub/seelog"

	"github.com/unirita/gocutoweb/config"
	"github.com/unirita/gocutoweb/testutil"
)

func initForTest() {
	config := `
<seelog type="sync" minlevel="trace">
    <outputs formatid="common">
        <console />
    </outputs>
    <formats>
        <format id="common" format="2015-04-01 12:34:56.789 [%LEV] %Msg%n"/>
    </formats>
</seelog>`
	logger, _ := seelog.LoggerFromConfigAsString(config)
	mutex = new(sync.Mutex)
	seelog.ReplaceLogger(logger)
	isValid = true
}

func TestInit_InvalidLogLevelError(t *testing.T) {
	config.Log.OutputLevel = "unknown"
	err := Init()
	if err == nil {
		t.Error("Error must be occured, but was not.")
	}
}

func TestTrace(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()

	c.Start()
	Trace("testmessage")
	output := c.Stop()

	expected := "2015-04-01 12:34:56.789 [TRC] testmessage\n"
	if output != expected {
		t.Errorf("Unexpected output message")
		t.Logf("Expected: %s", expected)
		t.Logf("Actual: %s", output)
	}
}

func TestTrace_Invalid(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()
	isValid = false

	c.Start()
	Trace("testmessage")
	output := c.Stop()

	expected := ""
	if output != expected {
		t.Errorf("Unexpected output message")
		t.Logf("Expected: %s", expected)
		t.Logf("Actual: %s", output)
	}
}

func TestDebug(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()

	c.Start()
	Debug("testmessage")
	output := c.Stop()

	expected := "2015-04-01 12:34:56.789 [DBG] testmessage\n"
	if output != expected {
		t.Errorf("Unexpected output message")
		t.Logf("Expected: %s", expected)
		t.Logf("Actual: %s", output)
	}
}

func TestDebug_Invalid(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()
	isValid = false

	c.Start()
	Debug("testmessage")
	output := c.Stop()

	expected := ""
	if output != expected {
		t.Errorf("Unexpected output message")
		t.Logf("Expected: %s", expected)
		t.Logf("Actual: %s", output)
	}
}

func TestInfo(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()

	c.Start()
	Info("testmessage")
	output := c.Stop()

	expected := "2015-04-01 12:34:56.789 [INF] testmessage\n"
	if output != expected {
		t.Errorf("Unexpected output message")
		t.Logf("Expected: %s", expected)
		t.Logf("Actual: %s", output)
	}
}

func TestInfo_Invalid(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()
	isValid = false

	c.Start()
	Info("testmessage")
	output := c.Stop()

	expected := ""
	if output != expected {
		t.Errorf("Unexpected output message")
		t.Logf("Expected: %s", expected)
		t.Logf("Actual: %s", output)
	}
}

func TestWarn(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()

	c.Start()
	Warn("testmessage")
	output := c.Stop()

	expected := "2015-04-01 12:34:56.789 [WRN] testmessage\n"
	if output != expected {
		t.Errorf("Unexpected output message")
		t.Logf("Expected: %s", expected)
		t.Logf("Actual: %s", output)
	}
}

func TestWarn_Invalid(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()
	isValid = false

	c.Start()
	Warn("testmessage")
	output := c.Stop()

	expected := ""
	if output != expected {
		t.Errorf("Unexpected output message")
		t.Logf("Expected: %s", expected)
		t.Logf("Actual: %s", output)
	}
}

func TestError(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()

	c.Start()
	Error("testmessage")
	output := c.Stop()

	expected := "2015-04-01 12:34:56.789 [ERR] testmessage\n"
	if output != expected {
		t.Errorf("Unexpected output message")
		t.Logf("Expected: %s", expected)
		t.Logf("Actual: %s", output)
	}
}

func TestError_Invalid(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()
	isValid = false

	c.Start()
	Error("testmessage")
	output := c.Stop()

	expected := ""
	if output != expected {
		t.Errorf("Unexpected output message")
		t.Logf("Expected: %s", expected)
		t.Logf("Actual: %s", output)
	}
}

func TestCritical(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()

	c.Start()
	Critical("testmessage")
	output := c.Stop()

	expected := "2015-04-01 12:34:56.789 [CRT] testmessage\n"
	if output != expected {
		t.Errorf("Unexpected output message")
		t.Logf("Expected: %s", expected)
		t.Logf("Actual: %s", output)
	}
}

func TestCritical_Invalid(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()
	isValid = false

	c.Start()
	Critical("testmessage")
	output := c.Stop()

	expected := ""
	if output != expected {
		t.Errorf("Unexpected output message")
		t.Logf("Expected: %s", expected)
		t.Logf("Actual: %s", output)
	}
}
