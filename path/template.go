package path

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/microhod/clone/repo"
)

// pathTemplate is a simple wrapper around text/template.Template
type pathTemplate struct {
	textTemplate *template.Template
}

// parsePathTemplate returns a pathTemplate from a string text template.
// In contrast to standard text templates, an empty string is not valid.
func parsePathTemplate(t string) (*pathTemplate, error) {
	if t == "" {
		return nil, fmt.Errorf("cannot create pathTemplate from an empty string")
	}

	textTemplate, err := template.New("PathTemplate").Parse(t)
	if err != nil {
		return nil, err
	}
	return &pathTemplate{textTemplate: textTemplate}, nil
}

// executeForRepo executes the template against the passed in repo and returns a string path
func (t *pathTemplate) executeForRepo(r repo.Repo) (string, error) {
	builder := new(strings.Builder)
	err := t.textTemplate.Execute(builder, r)
	if err != nil {
		return "", fmt.Errorf("failed to execute PathTemplate on repo: %s", err)
	}

	return builder.String(), nil
}
