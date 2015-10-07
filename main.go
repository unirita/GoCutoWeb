package main

import (
	"fmt"
	"net/http"
)

func main() {
	handler := setupHandler()

	if err := http.ListenAndServe(":8000", handler); err != nil {
		fmt.Println(err)
	}
}
