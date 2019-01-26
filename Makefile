format:
	docker run --rm -v "${PWD}":/app -w /app golang:1.11 go fmt ./src

run:
	docker run --rm -v "${PWD}":/app -w /app golang:1.11 go run src/main.go

test:
	docker run --rm -v "${PWD}":/app -w /app golang:1.11 go test -v ./src

benchmark:
	docker run --rm -v "${PWD}":/app -w /app golang:1.11 go test -v -bench=. ./src

check: format test
