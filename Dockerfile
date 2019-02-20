FROM golang:1.11

RUN go get -u golang.org/x/tools/cmd/goimports
RUN go get -u github.com/sergi/go-diff/...
