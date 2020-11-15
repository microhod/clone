package repo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	repoAPI = "https://api.github.com/repos"
)

// GetMainLanguageGithub gets the main language for the github repository
func GetMainLanguageGithub(r Repo) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/%s/languages", repoAPI, r.Owner, r.Repo))
	if err != nil {
		return "", fmt.Errorf("could not get repo languages from github: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not get repo language response body from github: %s", err)
	}

	var languages map[string]int
	err = json.Unmarshal(body, &languages)
	if err != nil {
		return "", fmt.Errorf("could not parse repo language response from github: %s", err)
	}

	return strings.ToLower(max(languages)), nil
}

func max(m map[string]int) string {
	max := 0
	key := ""
	for k, v := range m {
		if v > max {
			max = v
			key = k
		}
	}
	return key
}
