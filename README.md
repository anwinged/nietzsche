# Nietzsche

![Nietzsche](/share/Nietzsche.jpg)

Mustache template renderer. Cause Nietzsche write novels and has a great mustache.

## Build

Simply build with docker:

	$ make build-docker
	$ make build

## Usage

Render template

	$ nietzsche template.mustache data.json

Print template structure

	$ nietzsche -tree template.mustache

## Tests

	$ make test

Run tests in docker, so be sure you build environment with `make build-docker`.
