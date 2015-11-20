all:
	docker run --rm -v $(CURDIR):/go golang:1.5.1 \
		go build -o /go/build/waitforservices
