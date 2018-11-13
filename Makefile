all:
	docker run --rm -e CGO_ENABLED=0 -v $(CURDIR):/go golang:1.11.2 \
		go build -o /go/build/waitforservices
