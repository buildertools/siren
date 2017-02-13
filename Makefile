PWD := $(shell pwd)

prepare:
	@docker build -t buildertools/siren:build-tooling -f tooling.df .

update-deps:
	@docker run --rm \
          -v $(PWD):/go/src/github.com/buildertools/siren \
          -w /go/src/github.com/buildertools/siren \
          buildertools/siren:build-tooling trash -u
update-vendor:
	@docker run --rm \
          -v $(PWD):/go/src/github.com/buildertools/siren \
          -w /go/src/github.com/buildertools/siren \
          buildertools/siren:build-tooling trash

test:
	@docker run --rm \
	  -v $(PWD):/go/src/github.com/buildertools/siren \
	  -v $(PWD)/bin:/go/bin \
	  -v $(PWD)/pkg:/go/pkg \
	  -v $(PWD)/reports:/go/reports \
	  -w /go/src/github.com/buildertools/siren \
	  golang:1.7 \
	  go test -cover ./...
	  
