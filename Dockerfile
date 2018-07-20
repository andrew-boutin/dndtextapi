FROM golang:1.9

COPY . /go/src/github.com/andrew-boutin/dndtextapi
WORKDIR /go/src/github.com/andrew-boutin/dndtextapi

RUN go get -u github.com/kardianos/govendor
RUN govendor install +local
RUN go build
	
EXPOSE 8080