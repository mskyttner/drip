# README

Running drip in a container, with R and plumber installed.

When developing plumber APIs, it would be nice to be able to mount a local PLUMBER_DIR into a container for "live reload" changes to appear instantly at for example http://localhost:8000. 

## Usage

Use the Makefile targets:

		make build
		make test-default
		make up

## Issues / questions

When developing plumber files located locally in the PLUMBER_DIR, changes are expected to immediately drip (live deploy/reload) in the plumber server.

A few observations and questions:

- The default plumber example (gets installed with plumber R package) does not run pr$run() which drip expects?
- It seems drip needs to run in the present directory for the entrypoint.R used in the example? For details, try to change the `drip --dir` parameter with first changing present working directory.
- After running `make up` the curl command (see the Makefile) does not increment the counter in the example?


