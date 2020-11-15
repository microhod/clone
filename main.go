package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/microhod/clone/path"
	"github.com/microhod/clone/repo"
)

// Config describes the configuration for this application
type Config struct {
	DefaultHost      string            `json:"defaultHost"`
	PathTemplates    map[string]string `json:"pathTemplates"`
	DefaultProtocals map[string]string `json:"defaultProtocals"`
}

func main() {
	var err error

	l := flag.String("l", "", "optionally set the main language of the repo")
	flag.Parse()
	lang := *l
	if len(flag.Args()) < 1 {
		log.Println("please set the repository you want to clone")
		os.Exit(1)
	}

	var config *Config
	if configPath := os.Getenv("CLONE_CONFIG"); configPath != "" {
		config, err = parseConfig(configPath)
		if err != nil {
			log.Printf("error parsing config file '%s': %s", configPath, err)
		}
	}
	if config == nil {
		config = defaultConfig()
	}

	rp := repo.NewParser(config.DefaultHost, config.DefaultProtocals)
	repo, err := rp.Parse(flag.Args()[0])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if lang == "" {
		lang, err = repo.GetMainLanguage()
		if err != nil {
			log.Println(err)
		}
	}

	p := path.NewParser(config.PathTemplates)
	path, err := p.Parse(repo, lang)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Println(err)
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
			"go":      "~/go/src/${host}/${owner}/${repo}",
			"default": "~/src/${host}/${owner}/${repo}",
		},
		DefaultProtocals: map[string]string{
			"github.com": "git@",
			"default":    "https://",
		},
	}
}
