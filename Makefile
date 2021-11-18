BINARY := flextime.bin
VERSION := 2021-10-30
UNAME := $(shell uname -s)
SOURCES := $(wildcard flex/*.go cmd/*.go)
COMMIT_ID := $(shell git describe --tags --always)
BUILD_TIME := $(shell go run tool/rfc3339date.go)
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.BUILD_DATE=${BUILD_TIME} -X main.COMMIT_ID=${COMMIT_ID} -s -w ${DFLAG}"

ifeq ($(UNAME), Linux)
	DFLAG := -d
endif

.DEFAULT_GOAL: $(BINARY)

# Since we have build constraints, we should pass '.' (package) to build, not a list of go files
$(BINARY): $(SOURCES)
	cd cmd && env CGO_ENABLED=0 go build ${LDFLAGS} -o ../$@ .

.PHONY: install
install:
	cd cmd && env CGO_ENABLED=0 go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ]; then rm -f ${BINARY}; fi

#.PHONY: dbg
#dbg:
#	echo env CGO_ENABLED=0 go build ${LDFLAGS} -o $@ ${SOURCES}
