package define

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/unirita/gocutoweb/log"
)

// find jobnet json template and replace.
func ReplaceJobnetTemplate(path, jobnetName string, params []string) (string, error) {
	jobnetPath := filepath.Join(path, jobnetName+".json")
	log.Debug("Jobnet file path: ", jobnetPath)

	templateFile, err := os.Open(jobnetPath)
	if err != nil {
		return "", err
	}
	defer templateFile.Close()

	jsonBuf, err := ioutil.ReadAll(templateFile)
	if err != nil {
		return "", err
	}

	cacheName, err := replaceAndSave(string(jsonBuf), params)
	if err != nil {
		return "", err
	}
	return cacheName, nil
}

func replaceAndSave(jobnetJson string, params []string) (string, error) {
	jobnetJson = ExpandVariables(jobnetJson, params...)
	// save
	cacheName := time.Now().Format("20060102150405.000")
	os.MkdirAll(filepath.Join(os.TempDir(), "gocuto"), 0777)
	f, err := os.Create(filepath.Join(os.TempDir(), "gocuto", cacheName+".json"))
	if err != nil {
		return "", err
	}
	f.WriteString(jobnetJson)
	f.Close()

	return cacheName, nil
}
