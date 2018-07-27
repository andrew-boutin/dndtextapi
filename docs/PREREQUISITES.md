# Prerequisites

## Setup

Install [`docker`](https://www.docker.com/) and [`docker-compose`](https://docs.docker.com/compose/).

Set environment variable `DNDTEXTAPI_ENV` to either `int` or `prod`.

`int` will use the config file `config-int.yml` which is already set up. This uses authentication against a mock server int the  compose network which allows the integration tests to authenticate with the server.

`prod` will use the config file `config-prod.yml` which you will have to set up. Choose your `postgresql` configuration info. Right now authentication only works with Google. Set the following `accounts: https://www.googleapis.com` and `oauth2: https://accounts.google.com`. For the `id` and `secret` you'll have to set up a free [`Google Cloud Project`](https://console.cloud.google.com). This will give you a client id and secret. The callback URL in the Google cloud project config should be `http://localhost:8080/callback`.

## Development

Install [`Go`](https://golang.org/).

In order to get dependencies set up `Govendor`:

```bash
# Install
go get -u github.com/kardianos/govendor

# Set environement variable that Makefile utilizes
export GOVENDOR_PATH=<govendor binary location - possibly $GOPATH/bin>
```

In order to lint set up the following:

```bash
# Install
go get -u golang.org/x/lint/golint

# Set environement variable that Makefile utilizes
export GOLINT_PATH=<golint binary location - possibly $GOPATH/bin>
```