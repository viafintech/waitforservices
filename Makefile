all:
	docker run --rm -v $(CURDIR):/go golang:1.4.2 \
		go build -o /go/build/waitforservices
