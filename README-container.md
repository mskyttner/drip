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
- Reloads sometimes cause plumber to halt, while the drip process continues. Even if using "restart: always" on the container, won't help because drip is still fine. Could the exit status of plumber (if execution halts) be propagated somehow so that drip can restart it?

When this happens it seems there might be a port clash?

		server_1  | createTcpServer: address already in use
		server_1  | Error in initialize(...) : Failed to create server
		server_1  | Calls: <Anonymous> ... <Anonymous> -> startServer -> <Anonymous> -> initialize
		server_1  | Execution halted

