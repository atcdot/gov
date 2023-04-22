package dlsdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

const goVersionsURL = "https://go.dev/dl/?mode=json&include=all"

func GetVersions() ([]string, error) {
	resp := DlResp{}
	err := getJson(goVersionsURL, &resp)
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(resp))
	for _, r := range resp {
		isAvailableArch := false
		for _, file := range r.Files {
			if file.Arch == runtime.GOARCH && file.OS == runtime.GOOS {
				isAvailableArch = true
				break
			}
		}
		if r.Stable && isAvailableArch {
			versions = append(versions, r.Version)
		}
	}

	return versions, nil
}

func VersionsAvailable() DlResp {
	resp := DlResp{}
	err := getJson(goVersionsURL, &resp)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	rr := DlResp{}
	for _, r := range resp {
		isAvailableArch := false
		for _, file := range r.Files {
			if file.Arch == runtime.GOARCH && file.OS == runtime.GOOS {
				isAvailableArch = true
				break
			}
		}
		if r.Stable && isAvailableArch {
			rr = append(rr, r)
		}
	}

	return rr
}

type DlResp []Release

type Release struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []File `json:"files"`
}

type File struct {
	Filename       string `json:"filename"`
	OS             string `json:"os"`
	Arch           string `json:"arch"`
	Version        string `json:"version"`
	ChecksumSHA256 string `json:"sha256" datastore:",noindex"`
	Size           int64  `json:"size" datastore:",noindex"`
	Kind           string `json:"kind"` // "archive", "installer", "source"
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
