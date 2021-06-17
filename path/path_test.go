package path

import (
	"fmt"
	"os/user"
	"strings"
	"testing"

	"github.com/microhod/clone/repo"
)

func TestNewPathParserWithValidTemplatesReturnsPathParser(t *testing.T) {
	// arrange
	templates := map[string]string{
		"go": "{{.URL}}/{{.Scheme}}/{{.User}}/{{.Host}}/{{.Owner}}/{{.Repo}}",
	}

	// act
	parser, err := NewPathParser(templates)

	// assert
	if err != nil {
		t.Errorf("failed to create PathParser: %s", err)
	}
	if parser == nil {
		t.Errorf("expected parser not to be nil, but got nil")
	}
}

func TestNewPathParserWithNoTemplatesReturnsPathParser(t *testing.T) {
	// arrange
	templates := map[string]string{}

	// act
	parser, err := NewPathParser(templates)

	// assert
	if err != nil {
		t.Errorf("failed to create PathParser: %s", err)
	}
	if parser == nil {
		t.Errorf("expected parser not to be nil, but got nil")
	}
}

func TestNewPathParserWithInvalidTemplatesReturnsError(t *testing.T) {
	// arrange
	templates := map[string]string{
		"go": "{{}",
	}

	// act
	parser, err := NewPathParser(templates)

	// assert
	expectedErrorPrefix := "failed to parse PathTemplate from string '{{}'"
	if err == nil || !strings.HasPrefix(err.Error(), expectedErrorPrefix) {
		t.Errorf("expected non-nil error, starting with '%s': %s", expectedErrorPrefix, err)
	}
	if parser != nil {
		t.Errorf("expected PathTemplate to be nil, but got: %+v", parser)
	}
}

func TestParseReturnsCorrectString(t *testing.T) {
	// arrange
	templates := map[string]string{
		"go":      "{{.URL}}/{{.Scheme}}/{{.User}}/{{.Host}}/{{.Owner}}/{{.Repo}}",
		"default": "default/path",
	}

	// act
	parser, err := NewPathParser(templates)
	if err != nil {
		t.Errorf("failed to create PathParser: %s", err)
	}
	path, err := parser.Parse("go", repo.Repo{
		URL:    "repo.url",
		Scheme: "scheme",
		User:   "user",
		Host:   "host",
		Owner:  "owner",
		Repo:   "repo",
	})

	// assert
	if err != nil {
		t.Errorf("failed to parse path: %s", err)
	}
	expected := "repo.url/scheme/user/host/owner/repo"
	if path != expected {
		t.Errorf("expected '%s', but got '%s'", expected, path)
	}
}

func TestParseWithDifferentLanguageCasesReturnsCorrectString(t *testing.T) {
	// arrange
	templates := map[string]string{
		"Go":      "{{.URL}}/{{.Scheme}}/{{.User}}/{{.Host}}/{{.Owner}}/{{.Repo}}",
		"default": "default/path",
	}

	// act
	parser, err := NewPathParser(templates)
	if err != nil {
		t.Errorf("failed to create PathParser: %s", err)
	}
	path, err := parser.Parse("gO", repo.Repo{
		URL:    "repo.url",
		Scheme: "scheme",
		User:   "user",
		Host:   "host",
		Owner:  "owner",
		Repo:   "repo",
	})

	// assert
	if err != nil {
		t.Errorf("failed to parse path: %s", err)
	}
	expected := "repo.url/scheme/user/host/owner/repo"
	if path != expected {
		t.Errorf("expected '%s', but got '%s'", expected, path)
	}
}

func TestParseWithDefaultReturnsCorrectString(t *testing.T) {
	// arrange
	templates := map[string]string{
		"default": "{{.URL}}/{{.Scheme}}/{{.User}}/{{.Host}}/{{.Owner}}/{{.Repo}}",
	}

	// act
	parser, err := NewPathParser(templates)
	if err != nil {
		t.Errorf("failed to create PathParser: %s", err)
	}
	path, err := parser.Parse("go", repo.Repo{
		URL:    "repo.url",
		Scheme: "scheme",
		User:   "user",
		Host:   "host",
		Owner:  "owner",
		Repo:   "repo",
	})

	// assert
	if err != nil {
		t.Errorf("failed to parse path: %s", err)
	}
	expected := "repo.url/scheme/user/host/owner/repo"
	if path != expected {
		t.Errorf("expected '%s', but got '%s'", expected, path)
	}
}

func TestParseWithGlobalDefaultReturnsCorrectString(t *testing.T) {
	// arrange
	templates := map[string]string{}

	// act
	parser, err := NewPathParser(templates)
	if err != nil {
		t.Errorf("failed to create PathParser: %s", err)
	}
	path, err := parser.Parse("go", repo.Repo{
		URL:    "repo.url",
		Scheme: "scheme",
		User:   "user",
		Host:   "host",
		Owner:  "owner",
		Repo:   "repo",
	})

	// assert
	if err != nil {
		t.Errorf("failed to parse path: %s", err)
	}
	usr, _ := user.Current()
	expected := fmt.Sprintf("%s/src/host/owner/repo", usr.HomeDir)
	if path != expected {
		t.Errorf("expected '%s', but got '%s'", expected, path)
	}
}

func TestParseWithTildaReturnsCorrectStringWithUserDirReplaced(t *testing.T) {
	// arrange
	templates := map[string]string{
		"go": "~/{{.Host}}/~shouldnotbereplacedhere",
	}

	// act
	parser, err := NewPathParser(templates)
	if err != nil {
		t.Errorf("failed to create PathParser: %s", err)
	}
	path, err := parser.Parse("go", repo.Repo{Host: "host"})

	// assert
	if err != nil {
		t.Errorf("failed to parse path: %s", err)
	}
	usr, _ := user.Current()
	expected := fmt.Sprintf("%s/host/~shouldnotbereplacedhere", usr.HomeDir)
	if path != expected {
		t.Errorf("expected '%s', but got '%s'", expected, path)
	}
}

func TestParseWithInvalidTemplateReturnsError(t *testing.T) {
	// arrange
	templates := map[string]string{
		"go":      "{{.FieldDoesNotExist}}",
		"default": "default/path",
	}

	// act
	parser, err := NewPathParser(templates)
	if err != nil {
		t.Errorf("failed to create PathParser: %s", err)
	}
	path, err := parser.Parse("go", repo.Repo{
		URL:    "repo.url",
		Scheme: "scheme",
		User:   "user",
		Host:   "host",
		Owner:  "owner",
		Repo:   "repo",
	})

	// assert
	expectedErrorPrefix := "failed to execute PathTemplate:"
	if err == nil || !strings.HasPrefix(err.Error(), expectedErrorPrefix) {
		t.Errorf("expected non-nil error, starting with '%s': %s", expectedErrorPrefix, err)
	}
	if path != "" {
		t.Errorf("expected path to be empty, but got: %s", path)
	}
}
