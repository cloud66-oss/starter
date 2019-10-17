package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"errors"
)


type DownloadFile struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type TemplateDefinition struct {
	Version           string         `json:"version"`
	Dockerfiles       []DownloadFile `json:"dockerfiles"`
	ServiceYmls       []DownloadFile `json:"service-ymls"`
	DockerComposeYmls []DownloadFile `json:"docker-compose-ymls"`
	BundleManifest    []DownloadFile `json:"bundle-manifest-jsons"`
}


func Fetch(url string, mod *time.Time) (io.ReadCloser, error) {
	PrintlnL2("Downloading from %s", url)

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

func FetchJSON(url string, mod *time.Time, v interface{}) error {
	r, err := Fetch(url, mod)
	if err != nil {
		return err
	}

	defer r.Close()
	return json.NewDecoder(r).Decode(v)
}


func DownloadTemplates(tempDir string, td TemplateDefinition, templatePath string, flagBranch string) error {
	err := DownloadSingleFile(tempDir, DownloadFile{URL: strings.Replace(templatePath, "{{.branch}}", flagBranch, -1), Name: "templates.json"}, flagBranch)
	if err != nil {
		return err
	}

	for _, temp := range td.Dockerfiles {
		err := DownloadSingleFile(tempDir, temp, flagBranch)
		if err != nil {
			return err
		}
	}

	for _, temp := range td.ServiceYmls {
		err := DownloadSingleFile(tempDir, temp, flagBranch)
		if err != nil {
			return err
		}
	}

	for _, temp := range td.DockerComposeYmls {
		err := DownloadSingleFile(tempDir, temp, flagBranch)
		if err != nil {
			return err
		}
	}

	for _, temp := range td.BundleManifest {
		err := DownloadSingleFile(tempDir, temp, flagBranch)
		if err != nil {
			return err
		}
	}

	return nil
}

func DownloadSingleFile(tempDir string, temp DownloadFile, flagBranch string) error {
	r, err := Fetch(strings.Replace(temp.URL, "{{.branch}}", flagBranch, -1), nil)
	if err != nil {
		return err
	}
	defer r.Close()

	output, err := os.Create(filepath.Join(tempDir, temp.Name))
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, r)
	if err != nil {
		return err
	}

	return nil
}

