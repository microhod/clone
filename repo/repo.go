package repo

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	defaultScheme = "https://"
)

// Repo describes a repository by host (e.g. github.com), owner (e.g. microhod) and repo (e.g. )
type Repo struct {
	URL    string
	Scheme string
	User   string
	Host   string
	Owner  string
	Repo   string
}

// Parser holds config for parsing repos
type Parser struct {
	DefaultHost    string
	DefaultSchemes map[string]string
}

// NewParser creates a parser with the supplied default schemes
func NewParser(defaultHost string, defaultSchemes map[string]string) *Parser {
	return &Parser{
		DefaultHost:    defaultHost,
		DefaultSchemes: defaultSchemes,
	}
}

// Parse returns a repo object for a given string
func (p *Parser) Parse(repo string) (Repo, error) {
	r := parse(repo)

	if r.Owner == "" || r.Repo == "" {
		return Repo{}, fmt.Errorf("Could not parse owner and repo from input: %s", repo)
	}
	if r.Host == "" {
		r.Host = p.DefaultHost
	}
	if r.Scheme == "" {
		r.Scheme = p.DefaultSchemes[r.Host]
	}
	if r.Scheme == "" {
		r.Scheme = p.DefaultSchemes["default"]
	}
	if r.Scheme == "" {
		r.Scheme = defaultScheme
	}

	r.updateURL()
	return r, nil
}

func parse(repo string) Repo {
	r := Repo{URL: repo}
	if strings.HasPrefix(repo, "git@") {
		r.Scheme = "git@"
		repo = strings.Replace(repo, "git@", "", 1)
		repo = strings.Replace(repo, ":", "/", 1)
	} else {
		parts := strings.Split(repo, "://")
		if len(parts) > 1 {
			var re = regexp.MustCompile(`[^:\/\/]*:\/\/`)
			repo = re.ReplaceAllString(repo, "")
			r.Scheme = fmt.Sprintf("%s://", parts[0])
		}
	}

	parts := strings.Split(repo, "/")
	if len(parts) > 2 {
		r.Host = parts[0]
		r.Owner = parts[1]
		r.Repo = parts[2]
	} else if len(parts) > 1 {
		r.Owner = parts[0]
		r.Repo = parts[1]
	} else {
		r.Repo = parts[0]
	}

	return r
}

// GetMainLanguage returns the main language used in the repository
func (r *Repo) GetMainLanguage() (string, error) {
	if r.Host == "github.com" {
		return GetMainLanguageGithub(*r)
	}
	return "", nil
}

func (r *Repo) updateURL() {
	if r.Scheme == "git@" {
		r.URL = fmt.Sprintf("%s%s:%s/%s", r.Scheme, r.Host, r.Owner, r.Repo)
	} else {
		r.URL = fmt.Sprintf("%s%s/%s/%s", r.Scheme, r.Host, r.Owner, r.Repo)
	}
}
