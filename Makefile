#! make

build:
	# builds a test container with drip and R and plumber installed
	docker build --load -t siegerts/drip .

test-default:
	# seems not to work? because default example app at
	# /usr/local/lib/R/site-library/plumber/plumber/12-entrypoint
	# does not use pr$run()?
	docker run --rm siegerts/drip

up:
	# when adding pr$run() to the entrypoint.R
	# repeated requests to counter does not increment the counter?
	docker-compose up -d
	curl http://localhost:8000/counter
	curl http://localhost:8000/counter

