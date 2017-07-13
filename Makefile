
IMAGE=quay.io/ipedrazas
BINARY=backpressure
SHA1=$(shell git rev-parse HEAD | cut -c1-7)
VERSION=0.1.0-${SHA1}

clean:
	if [[ -f ${BINARY} ]] ; then rm ${BINARY} ; fi

test:
	go test -v 

build: clean
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BINARY} .
	docker build -t ${IMAGE}/${BINARY}:${VERSION} .

push:
	docker push ${IMAGE}/${BINARY}:${VERSION}


.PHONY: clean