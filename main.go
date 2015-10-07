package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/unirita/gocutoweb/config"
)

type arguments struct {
	configPath string
}

const defaultConfig = `web.ini`

func main() {
	args := fetchArgs()
	if args.configPath == "" {
		args.configPath = defaultConfig
	}

	if err := config.Load(args.configPath); err != nil {
		fmt.Println("Could not load config:", err)
	}

	listenHost := fmt.Sprintf(":%d", config.Server.ListenPort)
	if err := http.ListenAndServe(listenHost, setupHandler()); err != nil {
		fmt.Println("Could not start to listen:", err)
	}
}

func fetchArgs() *arguments {
	args := new(arguments)
	flag.StringVar(&args.configPath, "c", "", "Config file path")
	flag.Parse()
	return args
}
