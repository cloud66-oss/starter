package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cloud66-oss/starter/common"
)

func fetch(url string, mod *time.Time) (io.ReadCloser, error) {
	common.PrintlnL2("Downloading from %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if mod != nil {
		req.Header.Add("If-Modified-Since", mod.Format(http.TimeFormat))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if mod != nil && resp.StatusCode == 304 {
		return nil, errors.New("item moved")
	}
	if resp.StatusCode != 200 {
		err := fmt.Errorf("bad http status from %s: %v", url, resp.Status)
		return nil, err
	}
	if s := resp.Header.Get("Last-Modified"); mod != nil && s != "" {
		t, err := time.Parse(http.TimeFormat, s)
		if err == nil {
			*mod = t
		}
	}
	return resp.Body, nil
}

func fetchJSON(url string, mod *time.Time, v interface{}) error {
	r, err := fetch(url, mod)
	if err != nil {
		return err
	}

	defer r.Close()
	return json.NewDecoder(r).Decode(v)
}
