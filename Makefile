ifeq (,$(wildcard /.token))
	include .token
endif

install:
	go build && go install

uninstall:
	-rm ${GOBIN}/kubectl-plugin_completion

reinstall: uninstall install

release:
	goreleaser release --rm-dist

release-local:
	goreleaser release --snapshot --rm-dist

# export BASH_COMP_DEBUG_FILE=$(pwd)/p.sh
