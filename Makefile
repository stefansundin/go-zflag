.PHONY: lint test unittest zulutest clean

ifeq (, $(shell which golangci-lint 2> /dev/null))
	$(error Unable to locate golangci-lint! Ensure it is installed and available in PATH before re-running.)
endif

gotest =
ifeq (, $(shell which richgo 2> /dev/null))
	gotest = go test
else
	gotest = richgo test
endif

default: all

all: lint test clean

lint:
	@echo '********** LINT TEST **********'
	golangci-lint run

unittest:
	@echo '********** UNIT TEST **********'
	@$(gotest) -failfast -v -race -cover

zulutest:
	@echo '********** ZULU TEST **********'
	@set -e \
		&& test -d zulu || { git clone https://github.com/gowarden/zulu.git && ln -s ../../zflag zulu/zflag ; } \
		&& cd zulu \
		&& go mod edit -replace github.com/zulucmd/zflag=./zflag \
		&& $(gotest) -v ./...

test: unittest zulutest lint

clean:
	rm -rf zulu
