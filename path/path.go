package path

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/microhod/clone/repo"
)

const (
	defaultTemplateString = "~/src/{{.Host}}/{{.Owner}}/{{.Repo}}"
	defaultLanguageKey    = "default"
)

// PathParser contains the template mappig from language to PathTemplate
type PathParser struct {
	Templates map[string]*pathTemplate
}

// NewPathParser creates a PathParser from a map[string]string of languages to text templates
func NewPathParser(templates map[string]string) (*PathParser, error) {
	parserTemplates := map[string]*pathTemplate{}

	for language, templateString := range templates {
		var err error
		parserTemplates[strings.ToLower(language)], err = parsePathTemplate(templateString)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PathTemplate from string '%s': %s", templateString, err)
		}
	}

	return &PathParser{Templates: parserTemplates}, nil
}

// Parse computes the path using the relevant language template, and executing the template with the repo
func (p *PathParser) Parse(language string, r repo.Repo) (string, error) {
	template := p.Templates[strings.ToLower(language)]

	if template == nil {
		template = p.Templates[defaultLanguageKey]
	}
	if template == nil {
		var err error
		template, err = parsePathTemplate(defaultTemplateString)
		if err != nil {
			// this should never happen
			return "", fmt.Errorf("failed to parse the global default path template: %s", err)
		}
	}

	path, err := template.executeForRepo(r)
	if err != nil {
		return "", fmt.Errorf("failed to execute PathTemplate: %s", err)
	}

	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("could not get current os user to replace '~' in path: %s", err)
		}
		path = strings.Replace(path, "~", usr.HomeDir, 1)
	}

	return path, nil
}
