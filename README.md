# waitforservices

A small utility waiting for services linked to a Docker container being ready.

Without configuration, it finds all TCP services linked to a Docker container via their [environment variables](http://docs.docker.com/userguide/dockerlinks/#environment-variables) and concurrently and repeatedly tries to open a TCP connection to all of them.

When all connections are successful, it returns. If one or more services aren't ready within 60 seconds, it aborts and exits with status 1.

`waitforservices` also supports waiting for an HTTP request to `/` to return a response.

## Installation

First, install the utility into your image by adding this to your Dockerfile:

    RUN curl -LsS https://github.com/Barzahlen/waitforservices/releases/download/0.1/waitforservices \
            > /usr/local/bin/waitforservices && \
        chmod +x /usr/local/bin/waitforservices

There's also a Makefile in this repository if you want to build the binary yourself.

Then, during container startup, you can use the `waitforservices` command to wait for all services being ready.

## Usage

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

## License

waitforservices is licensed under [MIT](LICENSE).
