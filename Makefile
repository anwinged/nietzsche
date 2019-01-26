format:
	docker run --rm -v "${PWD}":/app -w /app golang:1.8 go fmt ./src

run:
	docker run --rm -v "${PWD}":/app -w /app golang:1.8 go run src/main.go

test:
	docker run --rm -v "${PWD}":/app -w /app golang:1.8 go test -v ./src