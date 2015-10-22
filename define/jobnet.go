package define

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/unirita/gocutoweb/log"
)

const REP_WEBAPI_PARM1 string = "$WAPARM1"
const REP_WEBAPI_PARM2 string = "$WAPARM2"

// find jobnet json template and replace.
func ReplaceJobnetTemplate(path, jobnetName string, params []string) (string, error) {
	jobnetPath := filepath.Join(path, jobnetName+".json")
	log.Debug("Jobnet file path: ", jobnetPath)

	template, err := os.Open(jobnetPath)
	if err != nil {
		return "", err
	}
	defer template.Close()

	network, err := Parse(template)
	if err != nil {
		return "", err
	}

	cacheName, err := network.replaceAndSave(params)
	if err != nil {
		return "", err
	}
	return cacheName, nil
}

func (n *Network) replaceAndSave(params []string) (string, error) {
	jobnetJson, err := n.Encode()
	if err != nil {
		return "", err
	}

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

func (n *Network) Encode() (string, error) {
	b, err := json.Marshal(n)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type Network struct {
	Flow string `json:"flow"`
	Jobs []Job  `json:"jobs"`
}

type Job struct {
	Name    string `json:"name"`
	Node    string `json:"node"`
	Port    int    `json:"port"`
	Path    string `json:"path"`
	Param   string `json:"param"`
	Env     string `json:"env"`
	Work    string `json:"work"`
	WRC     int    `json:"wrc"`
	WPtn    string `json:"wptn"`
	ERC     int    `json:"erc"`
	EPtn    string `json:"eptn"`
	Timeout int    `json:"timeout"`
	SNode   string `json:"snode"`
	SPort   int    `json:"sport"`
}

// Parse parses str as json format, and create Network object.
func Parse(reader io.Reader) (*Network, error) {
	decorder := json.NewDecoder(reader)

	network := new(Network)
	if err := decorder.Decode(network); err != nil {
		return nil, err
	}

	if err := network.DetectError(); err != nil {
		return nil, err
	}

	return network, nil
}

// DetectError detects error in Network object, and return it.
// If there is no error, DetectError returns nil.
func (n *Network) DetectError() error {
	for _, job := range n.Jobs {
		if job.Name == "" {
			return errors.New("Anonymous job detected.")
		}
	}
	return nil
}
