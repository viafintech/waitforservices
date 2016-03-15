# added

Multiple ports can be ignored with this fork by using e.g.

    waitforservices -ignoreport=8000,8001,8002

# waitforservices

A small utility waiting for services linked to a Docker container being ready.

When starting multiple Docker containers at once with containers depending on and [linking to other containers](http://docs.docker.com/userguide/dockerlinks/) (e.g. using [docker compose](https://github.com/docker/compose)), you might want to do some initalization in one container depending on a service in another container already running. E.g. a web application running database migrations on startup (for testing) might need a database service in a separate container to be running, but the database container might need a few seconds until it's started and ready for connections.

In your container startup script, waitforservices allows you to wait for other services to be ready by repeatedly trying to open a TCP connection to all linked services and blocking until it succeeds or times out.

We wrote a [blog post explaining why we built waitforservices how we use it](http://barzahlen.github.io/docker-waitforservices/).

## Installation

First, install the utility into your image by adding this to your Dockerfile:

    RUN curl --location --silent --show-error --fail \
            https://github.com/Barzahlen/waitforservices/releases/download/v0.3/waitforservices \
            > /usr/local/bin/waitforservices && \
        chmod +x /usr/local/bin/waitforservices

There's also a Makefile in this repository if you want to build the binary yourself.

Then, during container startup, you can use the `waitforservices` command to wait for all services being ready.

## Usage

Without configuration, it finds all TCP services linked to a Docker container via their [environment variables](http://docs.docker.com/userguide/dockerlinks/#environment-variables) and concurrently and repeatedly tries to open a TCP connection to all of them.

When all connections are successful, it returns. If one or more services aren't ready within a specified timeout (60 seconds by default), it aborts and exits with status 1.

`waitforservices` also supports waiting for an HTTP request to `/` to return a response.

    $ ./waitforservice -help
    Usage of ./waitforservices:
      -httpport=0: wait for an http request if target port is given port
      -ignoreport=0: don't wait for services on this port to be up
      -timeout=60: time to wait for all services to be up (seconds)

    Attempt to connect to all TCP services linked to a Docker container (found
    via their env vars) and wait for them to accept a TCP connection.

    When an <httpport> is specified, for services running on <httpport>, after
    a successful TCP connect, do an HTTP request and wait until it's done. This
    is useful for slow-starting services that only start up when they receive
    their first request.

    When timeout is over and TCP connect or HTTP request were unsucecssful, exit
    with status 1.

## Contributions

waitforservices' options are pretty limited at the moment (e.g. the `-httpport`  and `-ignoreport` parameters could support multiple ports), so we'd be happy if you create pull requests or report issues.

## License

waitforservices is licensed under [MIT](LICENSE).
