format:
	docker run --rm -v "${PWD}":/app -w /app golang:1.8 go fmt .

run:
	docker run --rm -v "${PWD}":/app -w /app golang:1.8 go run templater.go
