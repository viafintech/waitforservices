all:
	docker run --rm -e CGO_ENABLED=0 -v $(CURDIR):/go golang:1.5.1 \
		go build -o /go/build/waitforservices
