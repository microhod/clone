package path

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/microhod/clone/repo"
)

const (
	defaultTemplate = "~/src/${host}/${owner}/${repo}"
)

type Parser struct {
	Templates map[string]string
}

func NewParser(templates map[string]string) *Parser {
	return &Parser{Templates: templates}
}

// Parse parses the path using templates
func (p *Parser) Parse(r repo.Repo, lang string) (string, error) {
	t := p.Templates[lang]
	if t == "" {
		t = p.Templates["default"]
	}
	if t == "" {
		t = defaultTemplate
	}

	t = strings.ReplaceAll(t, "${host}", r.Host)
	t = strings.ReplaceAll(t, "${owner}", r.Owner)
	t = strings.ReplaceAll(t, "${repo}", r.Repo)

	if strings.HasPrefix(t, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("could not get current user: %s", err)
		}
		t = strings.Replace(t, "~", usr.HomeDir, 1)
	}

	return t, nil
}
