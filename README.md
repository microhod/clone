# Clone
This is a cli application to clone git repositories and automatically store them in a generated path according to templates and the main language of the repo.

## Examples
Here are some examples of clone commands, and their corresponding bash/git commands:

If no host is supplied, it will default to github
```
clone owner/repo -> 
    mkdir ~/src/github.com/owner/repo
    cd ~/src/github.com/owner/repo
    git clone git@github.com:owner/repo .
```
If no scheme is supplied, it will use the `defaultSchemes` from configuration, as described below
```
clone github.com/owner/repo -> 
    mkdir ~/src/github.com/owner/repo
    cd ~/src/github.com/owner/repo
    git clone git@github.com:owner/repo .
```
You can also specify a complete url to the repository
```
clone https://host.com/owner/repo -> 
    mkdir ~/src/host.com/owner/repo
    cd ~/src/host.com/owner/repo
    git clone https://host.com/owner/repo .
```

## Dependencies
[git](https://git-scm.com/downloads)

## Installation
Install simply by cloning the repository, running `go build` (this requires [go](https://golang.org/doc/install)) and placing the resulting executable somewhere on your PATH.

## Usage
```
clone -l <main_language> <repo_to_clone>
```
The main language can be autofilled if it is a github repository, using the github api.

On success, the command will output the absolute path to the repository that was cloned.

## Configuration
Configuration is provided as a json file, with deafaults as below:
```json
{
    "defaultHost": "github.com",
    "pathTemplates": {
        "go": "~/go/src/{{.Host}}/{{.Owner}}/{{.Repo}}",
        "default": "~/src/{{.Host}}/{{.Owner}}/{{.Repo}}"
    },
    "defaultSchemes": {
        "github.com": "git@",
        "default": "https://"
    }
}
```
This is provided to the application by setting the `CLONE_CONFIG` environment variable equal to the path to your config file.

### Path templates
Path templates describe the path in which to clone repositories, based on language (lower case). An example use case is golang code, which has to be stored under the `GOPATH`.

The `default` template is used if no other language template is available.

You can include any parts of the [Repo stuct](./repo/repo.go), which will be filled out at runtime (the templating uses [go text template](https://golang.org/pkg/text/template/) syntax)

You can also prefix your paths with `~`, which will be expanded to the current users home directory.

### Default Schemes
As mentioned in the examples, default schemes to use for `git clone` can be set per host. 

For instance in the default config, github repos default to the format `git@github.com:owner/repo`, whereas everything else defaults to `https`.
