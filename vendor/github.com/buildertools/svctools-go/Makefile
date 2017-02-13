PWD := $(shell pwd)

prepare:
	@docker build -t buildertools/svctools-go:build-tooling -f tooling.df .

update-deps:
	@docker run --rm -v $(PWD):/go/src/github.com/buildertools/svctools-go -w /go/src/github.com/buildertools/svctools-go buildertools/svctools-go:build-tooling trash -u
update-vendor:
	@docker run --rm -v $(PWD):/go/src/github.com/buildertools/svctools-go -w /go/src/github.com/buildertools/svctools-go buildertools/svctools-go:build-tooling trash

test:
	@docker run --rm \
	  -v $(PWD):/go/src/github.com/buildertools/svctools-go \
	  -v $(PWD)/bin:/go/bin \
	  -v $(PWD)/pkg:/go/pkg \
	  -v $(PWD)/reports:/go/reports \
	  -w /go/src/github.com/buildertools/svctools-go \
	  golang:1.7 \
	  go test -cover ./...
	  
build:
	@docker run --rm \
	  -v $(PWD):/go/src/github.com/buildertools/svctools-go \
	  -v $(PWD)/bin:/go/bin \
	  -v $(PWD)/pkg:/go/pkg \
	  -w /go/src/github.com/buildertools/svctools-go \
	  -e GOOS=darwin \
	  -e GOARCH=amd64 \
	  golang:1.7 \
	  go build -o bin/svctools
	  
