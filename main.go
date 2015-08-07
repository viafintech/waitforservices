package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Service struct {
	Name    string
	Address string
	Port    int
}

func (s Service) AddressAndPort() string {
	return fmt.Sprintf("%s:%d", s.Address, s.Port)
}

var timeout = flag.Int64("timeout", 60, "time to wait for all services to be up (seconds)")
var httpPort = flag.Int("httpport", 0, "wait for an http request if target port is given port")
var ignorePort = flag.Int("ignoreport", 0, "don't wait for services on this port to be up")

func main() {
	setupUsage()
	flag.Parse()

	services := loadServicesFromEnv()

	log.Printf("Waiting for %d services to be ready...", len(services))
	begin := time.Now()

	var wg sync.WaitGroup
	cancel := make(chan struct{})

	for _, service := range services {
		if service.Port == *ignorePort {
			continue
		}
		wg.Add(1)
		go func(service Service) {
			waitForTcpConn(service, cancel)
			if *httpPort == service.Port {
				waitForHttpRequest(service, cancel)
			}
			wg.Done()
		}(service)
	}

	timer := time.AfterFunc(time.Duration(*timeout)*time.Second, func() {
		close(cancel)
	})

	wg.Wait()

	// There's a race here that might result in assuming that a timeout happend
	// although none happend. It appears when the timer fires after the connection
	// succeeded, but before the check via Stop() below.
	// That shouldn't happen very often and the service was pretty short of timing out
	// anyway, so I guess that's ok for now.
	if !timer.Stop() {
		log.Printf("Error: One or more services timed out after %d second(s)", *timeout)
		os.Exit(1)
	}
	log.Printf("All services are up after %v!", time.Now().Sub(begin))
}

func setupUsage() {
	flag.CommandLine.Usage = func() {
		flag.Usage()
		fmt.Fprint(os.Stderr, `
Attempt to connect to all TCP services linked to a Docker container (found
via their env vars) and wait for them to accept a TCP connection.

When an <httpport> is specified, for services running on <httpport>, after
a successful TCP connect, do an HTTP request and wait until it's done. This
is useful for slow-starting services that only start up when they receive
their first request.

When timeout is over and TCP connect or HTTP request were unsucecssful, exit
with status 1.
`)
	}
}

func loadServicesFromEnv() []Service {
	services := make([]Service, 0)
	for _, line := range os.Environ() {
		keyAndValue := strings.SplitN(line, "=", 2)
		addrKey := keyAndValue[0]
		if strings.HasSuffix(addrKey, "_TCP_ADDR") {
			addr := os.Getenv(addrKey)
			name := addrKey[:len(addrKey)-9] // cut off "_TCP_ADDR"

			portKey := name + "_TCP_PORT"
			portStr := os.Getenv(portKey)
			port, err := strconv.Atoi(portStr)
			if err != nil {
				log.Printf("Failed to convert %v to int, value: '%v' - skipping service '%v'",
					portKey, portStr, name)
				continue
			}
			services = append(services, Service{Name: name, Address: addr, Port: port})
		}
	}
	return services
}

func waitForTcpConn(service Service, cancel <-chan struct{}) {
	var cancelled int32 = 0
	go func() {
		<-cancel
		atomic.StoreInt32(&cancelled, 1)
	}()

	var conn net.Conn
	err := errors.New("init")
	for err != nil {
		conn, err = net.DialTimeout("tcp", service.AddressAndPort(), 1*time.Second)

		if cancelled == 1 && err != nil {
			log.Printf("TCP: Service %v (%v) timed out. Last error: %v",
				service.Name, service.AddressAndPort(), err)
			return
		}

		time.Sleep(200 * time.Millisecond)
	}
	conn.Close()
	log.Printf("TCP: Service %v (%v) is up", service.Name, service.AddressAndPort())
}

func waitForHttpRequest(service Service, cancel <-chan struct{}) {
	url := url.URL{
		Scheme: "http",
		Host:   service.AddressAndPort(),
		Path:   "/",
	}

	err := errors.New("init")
	for err != nil {
		req, reqErr := http.NewRequest("GET", url.String(), nil)
		if reqErr != nil {
			log.Printf("Warning: Failed to create request for URL '%s' -  skipping service '%s'",
				url.String(), service.Name)
			return
		}

		tr := &http.Transport{}
		client := &http.Client{Transport: tr}
		c := make(chan error, 1)
		go func() {
			_, err := client.Do(req)
			c <- err
		}()

		select {
		case <-cancel:
			tr.CancelRequest(req)
			log.Printf("HTTP: Service %v (%v) timed out. Last error: %v",
				service.Name, service.AddressAndPort(), err)
			<-c
			return
		case err = <-c:
		}
		time.Sleep(200 * time.Millisecond)
	}

	log.Printf("HTTP: Service %v (%v) is up", service.Name, service.AddressAndPort())
}
