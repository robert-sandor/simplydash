#!/bin/sh

docker build \
  -f ./build/package/Dockerfile \
  -t simplydash:dev \
  .
