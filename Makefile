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
	${goexec} fmt

run:
	${goexec} run

test: format
	${goexec} test -v -cover

benchmark: format
	${goexec} test -bench=. -benchmem

coverage:
	${goexec} test -covermode=count -coverprofile=coverage.out
	${goexec} tool cover -html=coverage.out -o coverage.html

profile-cpu:
	${goexec} test -cpuprofile=profile-cpu.out -run=NONE -bench=.
	${goexec} tool pprof -text -nodecount=20 ./app.test profile-cpu.out

profile-mem:
	${goexec} test -memprofile=profile-mem.out -run=NONE -bench=.
	${goexec} tool pprof -text -nodecount=20 ./app.test profile-mem.out

clean:
	rm -rf ./.cache
	rm -f  ./coverage.{out,html}
	rm -f  ./profile-{cpu,mem}.out
	rm -f  ./app.test
