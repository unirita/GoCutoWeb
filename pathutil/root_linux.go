package pathutil

import (
	"fmt"
	"os"
)

var rootPath = getCutoRoot()

func getCutoRoot() string {
	d := os.Getenv("CUTOROOT")
	if len(d) == 0 {
		panic("Not setting environment argument $CUTOROOT")
	}
	return d
}

// Rootフォルダを取得する
func GetRootPath() string {
	return rootPath
}
