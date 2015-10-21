package define

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
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

	dynamicJobnetName, err := network.replaceAndSave(params)
	if err != nil {
		return "", err
	}
	return dynamicJobnetName, nil
}

func (n *Network) replaceAndSave(params []string) (string, error) {
	jobnetJson, err := n.Encode()
	if err != nil {
		return "", err
	}

	jobnetJson = ExpandVariables(jobnetJson, params...)
	// save
	dynamicJobnetName := time.Now().Format("20060102150405.000")
	os.MkdirAll(filepath.Join(os.TempDir(), "gocuto"), 0777)
	f, err := os.Create(filepath.Join(os.TempDir(), "gocuto", dynamicJobnetName+".json"))
	if err != nil {
		return "", err
	}
	f.WriteString(jobnetJson)
	f.Close()

	return dynamicJobnetName, nil
}

func (n *Network) Encode() (string, error) {
	b, err := json.Marshal(n)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

//Under Copied from realtime package's network.go.

// Number of columns
const columns = 14

// Indexes of columns
const (
	nameIdx = iota
	nodeIdx
	portIdx
	pathIdx
	paramIdx
	envIdx
	workIdx
	wrcIdx
	wptnIdx
	ercIdx
	eptnIdx
	timeoutIdx
	snodeIdx
	sportIdx
)

var jobex = make([][]string, 0)

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

// LoadJobex loads jobex csv which corresponds to name.
// LoadJobex returns empty jobex array if csv is not exists.
func LoadJobex(name string, nwkDir string) error {
	csvPath := searchJobexCsvFile(name, nwkDir)
	if csvPath == "" {
		return nil
	}

	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	jobex, err = loadJobexFromReader(file)
	return err
}

func searchJobexCsvFile(name string, nwkDir string) string {
	individualPath := filepath.Join(nwkDir, "realtime", name+".csv")
	defaultPath := filepath.Join(nwkDir, "realtime", "default.csv")

	if _, err := os.Stat(individualPath); !os.IsNotExist(err) {
		return individualPath
	}
	if _, err := os.Stat(defaultPath); !os.IsNotExist(err) {
		return defaultPath
	}

	return ""
}

// loadJobexFromReader reads reader as csv format, and create jobex data array.
func loadJobexFromReader(reader io.Reader) ([][]string, error) {
	r := csv.NewReader(reader)
	result := make([][]string, 0)
	isTitleRow := true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if !isTitleRow {
			result = append(result, record)
		}
		isTitleRow = false
	}
	if len(result) > 0 && len(result[0]) != columns {
		return nil, fmt.Errorf("Number of jobex csv columns[%d] must be %d.", len(result[0]), columns)
	}

	return result, nil
}

// Parse parses str as json format, and create Network object.
func Parse(reader io.Reader) (*Network, error) {
	decorder := json.NewDecoder(reader)

	network := new(Network)
	if err := decorder.Decode(network); err != nil {
		return nil, err
	}

	network.complementJobs()

	if err := network.DetectError(); err != nil {
		return nil, err
	}

	return network, nil
}

func (n *Network) complementJobs() {
	for _, record := range jobex {
		isExists := false
		for _, job := range n.Jobs {
			if record[nameIdx] == job.Name {
				isExists = true
				break
			}
		}

		if !isExists {
			newJob := Job{Name: record[nameIdx]}
			newJob.importJobex()
			n.Jobs = append(n.Jobs, newJob)
		}
	}
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

func (j *Job) importJobex() {
	for _, record := range jobex {
		if record[nameIdx] == j.Name {
			var err error
			j.Node = record[nodeIdx]
			j.Port, err = strconv.Atoi(record[portIdx])
			if err != nil {
				j.Port = 0
			}
			j.Path = record[pathIdx]
			j.Param = record[paramIdx]
			j.Env = record[envIdx]
			j.Work = record[workIdx]
			j.WRC, err = strconv.Atoi(record[wrcIdx])
			if err != nil {
				j.WRC = 0
			}
			j.WPtn = record[wptnIdx]
			j.ERC, err = strconv.Atoi(record[ercIdx])
			if err != nil {
				j.ERC = 0
			}
			j.EPtn = record[eptnIdx]
			j.Timeout, err = strconv.Atoi(record[timeoutIdx])
			if err != nil {
				j.Timeout = 0
			}
			j.SNode = record[snodeIdx]
			j.SPort, err = strconv.Atoi(record[sportIdx])
			if err != nil {
				j.SPort = 0
			}
		}
	}
}
