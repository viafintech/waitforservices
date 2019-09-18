# waitforservices ![Travis CI Status](https://travis-ci.org/Barzahlen/waitforservices.svg?branch=master)

A small utility waiting for services linked to a Docker container being ready.

## Installation

First, install the utility into your image by adding this to your Dockerfile:

    RUN curl --location --silent --show-error --fail \
            https://github.com/Barzahlen/waitforservices/releases/download/v0.6/waitforservices \
            > /usr/local/bin/waitforservices && \
        chmod +x /usr/local/bin/waitforservices

There's also a Makefile in this repository if you want to build the binary yourself.

Then, during container startup, you can use the `waitforservices` command to wait for all services being ready.

## Usage

Without configuration, it finds all TCP services specified by the environement variable declared like _\_HOST and _\_PORT (for example POSTGRES_HOST and POSTGRES_PORT) and concurrently and repeatedly tries to open a TCP connection to all of them.

When the _legacy_ option is specified, it finds all TCP services linked to a Docker container via their [environment variables](http://docs.docker.com/userguide/dockerlinks/#environment-variables)

When all connections are successful, it returns. If one or more services aren't ready within a specified timeout (60 seconds by default), it aborts and exits with status 1.

`waitforservices` also supports waiting for an HTTP request to `/` to return a response.

    $ ./waitforservice -help
    Usage of ./waitforservices:
      -httpport int  wait for an http request if target port is given port
      -ignoreport int don't wait for services on this port to be up
      -legacy use docker link enviroment variables
      -timeout int time to wait for all services to be up (seconds) (default 60)

    Attempt to connect to all TCP services linked  by the environement variable
    declared like _HOST and _PORT and wait for them to accept a TCP connection.

    When the _legacy_ option is specified, it finds all TCP services linked to
    a Docker container via their environment variables.

    When an <httpport> is specified, for services running on <httpport>, after
    a successful TCP connect, do an HTTP request and wait until it's done. This
    is useful for slow-starting services that only start up when they receive
    their first request.

    When timeout is over and TCP connect or HTTP request were unsucecssful, exit
    with status 1.

## Legacy support

When starting multiple Docker containers at once with containers depending on and [linking to other containers](http://docs.docker.com/userguide/dockerlinks/) (e.g. using [docker compose](https://github.com/docker/compose)), you might want to do some initalization in one container depending on a service in another container already running. E.g. a web application running database migrations on startup (for testing) might need a database service in a separate container to be running, but the database container might need a few seconds until it's started and ready for connections.

In your container startup script, waitforservices allows you to wait for other services to be ready by repeatedly trying to open a TCP connection to all linked services and blocking until it succeeds or times out.

We wrote a [blog post explaining why we built waitforservices how we use it](http://barzahlen.github.io/docker-waitforservices/).

## Contributions

waitforservices' options are pretty limited at the moment (e.g. the `-httpport` and `-ignoreport` parameters could support multiple ports), so we'd be happy if you create pull requests or report issues.

## License

waitforservices is licensed under [MIT](LICENSE).
