TEST ?= $(shell go list ./... | grep -v vendor)
NAME = "$(shell awk -F\" '/^const Name/ { print $$2; exit }' version.go)"
VERSION = "$(shell awk -F\" '/^const Version/ { print $$2; exit }' version.go)"

ifeq ("$(shell uname)","Darwin")
NCPU ?= $(shell sysctl hw.ncpu | cut -f2 -d' ')
else
NCPU ?= $(shell cat /proc/cpuinfo | grep processor | wc -l)
endif
TEST_OPTIONS=-timeout 30s -parallel $(NCPU)

default: test

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

depsdev: deps
	go get -u golang.org/x/lint/golint
	go get github.com/pierrre/gotestcover
	go get -u github.com/Songmu/goxz/cmd/goxz
	go get -u github.com/tcnksm/ghr
	go get -u github.com/Songmu/ghch/cmd/ghch

test:
	go test $(TEST) $(TESTARGS) $(TEST_OPTIONS)
	go test -race $(TEST) $(TESTARGS)

integration:
	go test -integration $(TEST) $(TESTARGS) $(TEST_OPTIONS)

lint:
	golint -set_exit_status $(TEST)

ci: depsdev test

clean:
	rm -rf ./builds
	mkdir ./builds

build: clean
	goxz -n $(NAME) -pv $(VERSION) -d ./builds -os=linux,darwin -arch=amd64 ./cmd/lolp

ghr:
	ghr -u pepabo v$(VERSION) builds

dist: build
	@test -z $(GITHUB_TOKEN) || $(MAKE) ghr

.PHONY: default dist test deps
