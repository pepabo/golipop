TEST ?= $(shell go list ./... | grep -v vendor)
NAME = "$(shell awk -F\" '/^const Name/ { print $$2; exit }' version.go)"
VERSION = "$(shell awk -F\" '/^const Version/ { print $$2; exit }' version.go)"

ifeq ("$(shell uname)","Darwin")
NCPU ?= $(shell sysctl hw.ncpu | cut -f2 -d' ')
CMD_INSTALL_GORELEASER="brew install goreleaser"
else
NCPU ?= $(shell cat /proc/cpuinfo | grep processor | wc -l)
CMD_INSTALL_GORELEASER="curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh"
endif
TEST_OPTIONS=-timeout 30s -parallel $(NCPU)

default: test

install_goreleaser:
	@which goreleaser || eval $(CMD_INSTALL_GORELEASER)

depsdev: export GO111MODULE=off
depsdev: install_goreleaser
	go get -u golang.org/x/lint/golint
	go get github.com/pierrre/gotestcover
	go get -u github.com/Songmu/ghch/cmd/ghch

test: export GO111MODULE=on
test:
	go test $(TEST) $(TESTARGS) $(TEST_OPTIONS)
	go test -race $(TEST) $(TESTARGS)

integration:
	go test -integration $(TEST) $(TESTARGS) $(TEST_OPTIONS)

lint:
	golint -set_exit_status $(TEST)

ci: test

build: install_goreleaser
	goreleaser --skip-publish --snapshot --rm-dist

dist: install_goreleaser
	goreleaser --rm-dist

.PHONY: default dist test
