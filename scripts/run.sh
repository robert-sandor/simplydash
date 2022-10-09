#!/bin/sh

docker run --rm -p 8080:8080 \
		-v /tmp/simplydash/config:/app/config \
		-v /tmp/simplydash/icons:/app/icons \
		--name simplydash-dev \
		simplydash
