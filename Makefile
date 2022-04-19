ifeq (,$(wildcard /.token))
	include .token
endif

install:
	go build && go install

release:
	goreleaser release --rm-dist

release-local:
	goreleaser release --snapshot --rm-dist
