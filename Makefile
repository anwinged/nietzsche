# -----------------------------------------------
# VARIABLES
# -----------------------------------------------

app_name := nietzsche
app_bin := ${app_name}
app_bin_test := ${app_bin}.test

build_dir := /app/build

uid := $(shell id -u)
gid := $(shell id -g)

image := ${app_name}-golang

goimage := docker run \
	--rm \
	--tty \
	--init \
	--user ${uid}:${gid} \
	--volume "${PWD}":/app \
	--env GOCACHE="/app/.cache" \
	--workdir /app \
	${image}

goexec := ${goimage} go

# -----------------------------------------------
# TARGETS
# -----------------------------------------------

.PHONY: build-docker
build-docker:
	docker build -t ${image} .

.PHONY: prepare-build-dir
prepare-build-dir:
	${goimage} mkdir -p ${build_dir}

.PHONY: format
format:
	${goimage} goimports -w .

.PHONY: build
build: format
	${goexec} build -o ${build_dir}/${app_bin}

.PHONY: test
test: format
	${goimage} gotest -v -cover

.PHONY: benchmark
benchmark: format
	${goexec} test -bench=. -benchmem

.PHONY: coverage
coverage: prepare-build-dir
	${goexec} test -covermode=count -coverprofile=${build_dir}/coverage.out
	${goexec} tool cover -html=${build_dir}/coverage.out -o ${build_dir}/coverage.html

.PHONY: profile-cpu
profile-cpu: prepare-build-dir
	${goexec} test -cpuprofile=${build_dir}/profile-cpu.out -run=NONE -bench=. -o ${build_dir}/${app_bin_test}
	${goexec} tool pprof -text -nodecount=20 ${build_dir}/${app_bin_test} ${build_dir}/profile-cpu.out

.PHONY: profile-mem
profile-mem: prepare-build-dir
	${goexec} test -memprofile=${build_dir}/profile-mem.out -run=NONE -bench=. -o ${build_dir}/${app_bin_test}
	${goexec} tool pprof -text -nodecount=20 ${build_dir}/${app_bin_test} ${build_dir}/profile-mem.out

clean:
	rm -rf ./.cache
	rm -rf ./build
