package path

import (
	"strings"
	"testing"

	"github.com/microhod/clone/repo"
)

func TestParsePathTemplateWithValidTemplateReturnsPathTemplate(t *testing.T) {
	// arrange
	var templateString = "some/path/{{.Some}}/{{.Value}}"

	// act
	template, err := parsePathTemplate(templateString)

	// assert
	if err != nil {
		t.Errorf("failed to parse path template: %s", err)
	}
	if template == nil {
		t.Errorf("expected PathTemplate not to be nil, but got nil")
	}
}

func TestParsePathTemplateWithInvalidTemplateReturnsError(t *testing.T) {
	// arrange
	var templateString = "{{}"

	// act
	template, err := parsePathTemplate(templateString)

	// assert
	if err == nil {
		t.Errorf("expected error not to be nil, but got nil")
	}
	if template != nil {
		t.Errorf("expected PathTemplate to be nil, but got: %+v", template)
	}
}

func TestParsePathTemplateWithEmptyTemplateReturnsError(t *testing.T) {
	// arrange
	var templateString = ""

	// act
	template, err := parsePathTemplate(templateString)

	// assert
	expectedErrorPrefix := "cannot create pathTemplate from an empty string"
	if err == nil || !strings.HasPrefix(err.Error(), expectedErrorPrefix) {
		t.Errorf("expected non-nil error, starting with '%s': %s", expectedErrorPrefix, err)
	}
	if template != nil {
		t.Errorf("expected PathTemplate to be nil, but got: %+v", template)
	}
}

func TestExecuteForRepoWithValidTemplateReturnsCorrectString(t *testing.T) {
	// arrange
	templateString := "{{.URL}}/{{.Scheme}}/{{.User}}/{{.Host}}/{{.Owner}}/{{.Repo}}"

	// act
	template, err := parsePathTemplate(templateString)
	if err != nil {
		t.Errorf("failed to parse path template: %s", err)
	}
	str, err := template.executeForRepo(repo.Repo{
		URL:    "repo.url",
		Scheme: "scheme",
		User:   "user",
		Host:   "host",
		Owner:  "owner",
		Repo:   "repo",
	})

	// assert
	if err != nil {
		t.Errorf("failed to execute template: %s", err)
	}
	expected := "repo.url/scheme/user/host/owner/repo"
	if str != expected {
		t.Errorf("expected '%s', but got '%s'", expected, str)
	}
}

func TestExecuteForRepoWithEmptyRepoReturnsCorrectString(t *testing.T) {
	// arrange
	templateString := "{{.URL}}/{{.Scheme}}/{{.User}}/{{.Host}}/{{.Owner}}/{{.Repo}}"

	// act
	template, err := parsePathTemplate(templateString)
	if err != nil {
		t.Errorf("failed to parse path template: %s", err)
	}
	str, err := template.executeForRepo(repo.Repo{})

	// assert
	if err != nil {
		t.Errorf("failed to execute template: %s", err)
	}
	expected := "/////"
	if str != expected {
		t.Errorf("expected '%s', but got '%s'", expected, str)
	}
}

func TestExecuteForRepoWithInvalidTemplateReturnsError(t *testing.T) {
	// arrange
	templateString := "{{.FieldThatDoesNotExist}}"

	// act
	template, err := parsePathTemplate(templateString)
	if err != nil {
		t.Errorf("failed to parse path template: %s", err)
	}
	str, err := template.executeForRepo(repo.Repo{
		URL:    "repo.url",
		Scheme: "scheme",
		User:   "user",
		Host:   "host",
		Owner:  "owner",
		Repo:   "repo",
	})

	// assert
	expectedErrorPrefix := "failed to execute PathTemplate on repo:"
	if err == nil || !strings.HasPrefix(err.Error(), expectedErrorPrefix) {
		t.Errorf("expected non-nil error, starting with '%s': %s", expectedErrorPrefix, err)
	}
	if str != "" {
		t.Errorf("expected empty string, but got '%s'", str)
	}
}
