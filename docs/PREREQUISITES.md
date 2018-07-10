# Prerequisites

Install `go`, `docker`, and `docker-compose`.

In order to `make lint` set up the following:

    go get -u golang.org/x/lint/golint

    export GOLINT_PATH=<golint binary location - possibly $GOPATH/bin>