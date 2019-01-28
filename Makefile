uid = $(shell id -u)
goimage = docker run --rm --tty --user "${uid}":"${uid}" --volume "${PWD}":/app --workdir /app golang:1.11
goexec = ${goimage} go

format:
	${goexec} fmt ./src

run:
	${goexec} run src/main.go

test: format
	${goexec} test -v ./src

benchmark: format
	${goexec} test -v -bench=. ./src
