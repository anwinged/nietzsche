FROM golang:1.12.1

RUN go get -u golang.org/x/tools/cmd/goimports
RUN go get -u github.com/sergi/go-diff/...
RUN go get -u gopkg.in/yaml.v2
