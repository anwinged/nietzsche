# Nietzsche

![Nietzsche](/share/Nietzsche.jpg)

Mustache template renderer. Cause Nietzsche wrote novels and has a great mustache.

## Build

Simply build with docker:

	$ make build-docker
	$ make build

After check `build` subdirectory.

## Usage

Render template to standard output:

	$ nietzsche template.mustache data.json

Print template structure:

	$ nietzsche -tree template.mustache

## Tests

	$ make test

Run tests in docker, so be sure you build environment with `make build-docker`.
