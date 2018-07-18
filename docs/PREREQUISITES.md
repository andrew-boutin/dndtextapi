# Prerequisites

## Running

Install [`docker`](https://www.docker.com/) and [`docker-compose`](https://docs.docker.com/compose/).

### Setup Google Authentication

For authentication to work, you'll need to set up a [`Google Cloud Project`](https://console.cloud.google.com) that can be used to utilize Google for authentication. You'll need your Client ID and Client secret once you set up your project.

Create file `client.json` at the root of the repo that looks like:

```json
{
    "id":"google-project-client-id-here",
    "secret":"google-project-client-secret-here"
}
```

## Development

Install [`Go`](https://golang.org/).

In order to get dependencies set up Govendor:

```bash
# Install
go get -u github.com/kardianos/govendor

# Set environement variable that Makefile utilizes
export GOLINT_PATH=<golint binary location - possibly $GOPATH/bin>
```

In order to lint set up the following:

```bash
# Install
go get -u golang.org/x/lint/golint

# Set environement variable that Makefile utilizes
export GOLINT_PATH=<golint binary location - possibly $GOPATH/bin>
```