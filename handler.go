package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

const (
	methodGet    = "GET"
	methodPost   = "POST"
	methodPut    = "PUT"
	methodDelete = "DELETE"
)

func setupHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/caches/{name:[0-9]{17}}", showJSONCache).Methods(methodGet)
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
