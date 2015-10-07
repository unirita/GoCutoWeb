package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestShowJSONCache(t *testing.T) {
	server := httptest.NewServer(setupHandler())
	defer server.Close()

	cacheDir := filepath.Join(os.TempDir(), "gocuto")
	os.MkdirAll(cacheDir, 0777)
	cacheFile := filepath.Join(cacheDir, "20151007123456789.json")
	f, err := os.Create(cacheFile)
	if err != nil {
		t.Fatalf("File create error: %s", err)
	}
	f.WriteString("abcd")
	f.Close()
	defer os.Remove(cacheFile)

	output := testGetMessages(t, server.URL+"/caches/20151007123456789")
	if output != "abcd" {
		t.Errorf("output => %s, want %s", output, "abcd")
	}
}

func testGetMessages(t *testing.T, requestURL string) string {
	res, err := http.Get(requestURL)
	if err != nil {
		t.Fatalf("Error occured: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Logf("URL[%s]", requestURL)
		t.Errorf("res.StatusCode => %d, want %d", res.StatusCode, http.StatusOK)
	}

	resMsg, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Responce read error occured: %s", err)
	}

	return string(resMsg)
}
