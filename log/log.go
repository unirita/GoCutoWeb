// Package log provides logger interface.
package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cihub/seelog"

	"github.com/unirita/gocutoweb/config"
)

const logFileName = "web.log"

var isValid = false
var mutex *sync.Mutex

// Init initializes logger.
func Init() error {
	logfile := filepath.Join(config.Log.OutputDir, logFileName)
	if err := makeFileIfNotExist(logfile); err != nil {
		return err
	}

	logconf := generateConfigString(logfile)
	logger, err := seelog.LoggerFromConfigAsString(logconf)
	if err != nil {
		return err
	}

	seelog.ReplaceLogger(logger)
	mutex = new(sync.Mutex)
	isValid = true

	return nil
}

// Trace outputs trace level message.
func Trace(msg ...interface{}) {
	if !isValid {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	seelog.Trace(msg...)
}

// Debug outputs debug level message.
func Debug(msg ...interface{}) {
	if !isValid {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	seelog.Debug(msg...)
}

// Info outputs info level message.
func Info(msg ...interface{}) {
	if !isValid {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	seelog.Info(msg...)
}

// Warn outputs warn level message.
func Warn(msg ...interface{}) {
	if !isValid {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	seelog.Warn(msg...)
}

// Error outputs error level message.
func Error(msg ...interface{}) {
	if !isValid {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	seelog.Error(msg...)
}

// Critical outputs critical level message.
func Critical(msg ...interface{}) {
	if !isValid {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	seelog.Critical(msg...)
}

func makeFileIfNotExist(logfile string) error {
	if _, err := os.Stat(logfile); !os.IsNotExist(err) {
		// ファイルが存在する場合は何もしない。
		// os.IsExistはerr=nilのときfalseを返すため、os.IsNotExistで判定している。
		return nil
	}

	file, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	file.Close()
	return nil
}

func generateConfigString(logfile string) string {
	format := `
<seelog type="sync" minlevel="%s">
    <outputs formatid="common">
        <rollingfile type="size" filename="%s" maxsize="%d" maxrolls="%d" />
    </outputs>
    <formats>
        <format id="common" format="%%Date(2006-01-02 15:04:05.000) [%d] [%%LEV] %%Msg%%n"/>
    </formats>
</seelog>`

	return fmt.Sprintf(
		format,
		config.Log.OutputLevel,
		logfile,
		config.Log.MaxSizeKB*1024,
		config.Log.MaxGeneration-1,
		os.Getpid())
}
