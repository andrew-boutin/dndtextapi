FROM golang:1.9

# TODO: Go back to the simple version
ARG app_env
ENV APP_ENV $app_env

COPY . /go/src/github.com/andrew-boutin/dndtextapi
WORKDIR /go/src/github.com/andrew-boutin/dndtextapi

RUN go get ./
RUN go build

CMD if [ ${APP_ENV} = production ]; \
	then \
	app; \
	else \
	go get github.com/pilu/fresh && \
	fresh; \
	fi
	
EXPOSE 8080