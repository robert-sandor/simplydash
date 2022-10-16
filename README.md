# simplydash

Yet another simple, configurable and fast dashboard for your homelab or homepage.

## Table of contents
- [Usage](#usage)
- [Description](#description)
- [Features](#features)
- [Roadmap](#roadmap)
- [Acknowledgements](#acknowledgements)

## Usage

Deploy simplydash using the method of your choice below.
Use the samples provided in [deploy/samples](deploy/samples) or create 
your own by following the [Configuration](docs/configuration.md) section.

### Docker compose (recommended)

Docker compose is the recommended way to deploy simplydash. 
Here's an example, which can also be found in [deploy/docker/docker-compose.yml](deploy/docker/docker-compose.yml)

```yaml
version: "3.8"

services:

  simplydash:
    container_name: simplydash
    hostname: simplydash
    image: ghcr.io/robert-sandor/simplydash:dev
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - /tmp/simplydash/config:/app/config
      - /tmp/simplydash/icons:/app/icons
```

### Docker 

You can also run simplydash using just docker, if using docker compose is not wanted.

```shell
docker run -d -p 8080:8080 \
		-v /tmp/simplydash/config:/app/config \
		-v /tmp/simplydash/icons:/app/icons \
		--name simplydash-dev \
		simplydash
```

## Description

This project aims to a a simple and fast web dashboard, that can be used as your
homelab's main page, or as just a homepage for your browser. Self hosting this is
easy using Docker.

### Why create yet another dashboard

The project started from a want to have a dashboard that combines the features of others
that are already available and that I've tried in the past.

- [Flame]() 
  - Although Flame is a very minimalistic and elegant dashboard with support for 
docker labels, it lacks support for categories and an easy way to manage custom links 
from a script
- [Homarr]()
  - While very feature rich it is a great dashboard overall, but it's too complex for 
my needs, and also hard to automate
- [Homer]()
  - Another great dashboard, close to what Simplydash aims to be, but i couldn't get 
used to the look, and it also lacks docker or k8s support

## Features

- [x] YAML configuration file, with hot reload
- [x] Links from files, with hot reload
- [x] Docker / Docker Compose deployments
- [x] Proper dark / light mode support
- [x] Customizable theme
- [x] Cache external icons

## Roadmap

- [ ] Support multiple profiles
- [ ] Support links from Docker labels
- [ ] K8s deployment using Helm charts
- [ ] Support links from K8s labels
- [ ] Search box in UI
- [ ] Logo

## Acknowledgements

- [walkxcode/Dashboard-Icons](https://github.com/walkxcode/Dashboard-Icons) - The default icons are used from this repository
- [gin-gonic/gin](https://github.com/gin-gonic/gin) - Backend framework
- [https://github.com/sveltejs/svelte](https://github.com/sveltejs/svelte) - UI framework
- [materialdesignicons](https://materialdesignicons.com/) - Some of the UI icons