package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/microhod/clone/path"
	"github.com/microhod/clone/repo"
)

// Config describes the configuration for this application
type Config struct {
	DefaultHost    string            `json:"defaultHost"`
	PathTemplates  map[string]string `json:"pathTemplates"`
	DefaultSchemes map[string]string `json:"defaultSchemes"`
}

func main() {
	var err error

	l := flag.String("l", "", "optionally set the main language of the repo")
	flag.Parse()
	lang := *l
	if len(flag.Args()) < 1 {
		println("please set the repository you want to clone")
		os.Exit(1)
	}

	var config *Config
	if configPath := os.Getenv("CLONE_CONFIG"); configPath != "" {
		config, err = parseConfig(configPath)
		if err != nil {
			fmt.Printf("error parsing config file '%s': %s", configPath, err)
		}
	}
	if config == nil {
		config = defaultConfig()
	}

	rp := repo.NewParser(config.DefaultHost, config.DefaultSchemes)
	repo, err := rp.Parse(flag.Args()[0])
	if err != nil {
		println(err)
		os.Exit(1)
	}

	if lang == "" {
		lang, err = repo.GetMainLanguage()
		if err != nil {
			println(err.Error())
		}
	}

	p, err := path.NewPathParser(config.PathTemplates)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	path, err := p.Parse(lang, repo)
	if err != nil {
		println(err)
		os.Exit(1)
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		println(err)
		os.Exit(1)
	}

	out, err := exec.Command("git", "clone", repo.URL, path).CombinedOutput()
	if err != nil {
		fmt.Printf("Git ERROR:\n%s\n", string(out))
		os.Exit(1)
	}

	fmt.Println(path)
}

func parseConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var config *Config
	err = json.Unmarshal(contents, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func defaultConfig() *Config {
	return &Config{
		DefaultHost: "github.com",
		PathTemplates: map[string]string{
			"go":      "~/go/src/{{.Host}}/{{.Owner}}/{{.Repo}}",
			"default": "~/src/{{.Host}}/{{.Owner}}/{{.Repo}}",
		},
		DefaultSchemes: map[string]string{
			"github.com": "git@",
			"default":    "https://",
		},
	}
}
