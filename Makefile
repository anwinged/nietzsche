goexec = docker run --rm -v "${PWD}":/app -w /app golang:1.11 go

format:
	${goexec} fmt ./src

run:
	${goexec} run src/main.go

test: format
	${goexec} test -v ./src

benchmark: format
	${goexec} test -v -bench=. ./src
