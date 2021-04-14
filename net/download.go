package net

import (
	"io/fs"
	"io/ioutil"
	"net/http"
)

func DownloadBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func DownloadFile(url string, path string, perm fs.FileMode) ([]byte, error) {
	buf, err := DownloadBytes(url)
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(path, buf, perm); err != nil {
		return nil, err
	}
	return buf, nil
}
