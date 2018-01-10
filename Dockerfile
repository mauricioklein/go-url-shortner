FROM golang

ENV GOPATH /go

# Install govendor
RUN go get -u github.com/kardianos/govendor

ADD . /go/src/github.com/mauricioklein/go-url-shortner
WORKDIR /go/src/github.com/mauricioklein/go-url-shortner

# Download project dependencies
RUN govendor sync

# Build the binary
RUN go build -i cmd/go-url-shortner/main.go

# Expose port 8080
EXPOSE 8080

ENTRYPOINT ./main
