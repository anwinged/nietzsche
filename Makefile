uid = $(shell id -u)
gid = $(shell id -g)
image = golang:1.11

goimage = docker run \
	--rm \
	--tty \
	--init \
	--user ${uid}:${gid} \
	--volume "${PWD}":/app \
	--env GOCACHE=".cache" \
	--workdir /app \
	${image}

goexec = ${goimage} go

format:
	${goexec} fmt ./src

run:
	${goexec} run ./src

test: format
	${goexec} test -v -cover ./src

benchmark: format
	${goexec} test -v -bench=. ./src
