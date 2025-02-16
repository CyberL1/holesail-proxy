package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func GetLatestCliVersion() (*GithubRelease, error) {
	resp, err := http.Get(GithubReleaseUrl)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	release := &GithubRelease{}
	err = json.Unmarshal(body, release)
	if err != nil {
		return nil, err
	}
	return release, nil
}
