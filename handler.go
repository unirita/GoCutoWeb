package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/unirita/gocutoweb/config"
	"github.com/unirita/gocutoweb/define"
	"github.com/unirita/gocutoweb/log"
	"github.com/unirita/gocutoweb/realtimeutil"
)

const (
	methodGet    = "GET"
	methodPost   = "POST"
	methodPut    = "PUT"
	methodDelete = "DELETE"
)

const realtimeErrorResult = `{
    "status":2,
    "message":"%s",
    "pid":0,
	"network":{"instance":0,"name":""}
}`

func setupHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc(`/caches/{name:\d{14}\.\d{3}}`, showJSONCache).Methods(methodGet)
	router.HandleFunc(`/notice`, noticeJobnet).Methods(methodPost)
	return router
}

func showJSONCache(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	name := params["name"]

	cachePath := filepath.Join(os.TempDir(), "gocuto", name+".json")

	file, err := os.Open(cachePath)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Write(buf[:n])
	}
}

func noticeJobnet(writer http.ResponseWriter, request *http.Request) {
	jobnetwork := request.FormValue("jobnetwork")
	params := getFormParams(request)

	log.Info("Receive trigger jobnetwork[", jobnetwork, "] params", params)
	// create jobnet json-file from template.
	dynamicJobnetName, err := define.ReplaceJobnetTemplate(config.Jobnet.JobnetDir, jobnetwork, params)
	if err != nil {
		log.Warn("Jobnetwork [", jobnetwork, "] not found. Reason: ", err)
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// execute realtime utility.
	url := "http://127.0.0.1:" + strconv.Itoa(config.Server.ListenPort) + "/caches/" + dynamicJobnetName
	c := realtimeutil.NewCommand(dynamicJobnetName, url)
	if err = c.Run(); err != nil {
		// realtime utility execute error.
		c.Result = fmt.Sprintf(realtimeErrorResult, err.Error())
	}
	//response
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Write([]byte(c.Result))
}

func getFormParams(req *http.Request) []string {
	args := make([]string, 0)
	for i := 1; ; i++ {
		f := req.FormValue(fmt.Sprintf("param%d", i))
		if len(f) == 0 {
			break
		}
		args = append(args, f)
	}
	return args
}
