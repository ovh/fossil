BUILD_DIR=build
PREFIX = /usr/local

CC=go build
GITHASH=$(shell git rev-parse HEAD)
DFLAGS=-race
CFLAGS=-ldflags "-X github.com/runabove/fossil/cmd.githash=$(GITHASH)"
CROSS=GOOS=linux GOARCH=amd64


rwildcard=$(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) $(filter $(subst *,%,$2),$d))
VPATH= $(BUILD_DIR)

.SECONDEXPANSION:

build: fossil.go $$(call rwildcard, ./cmd, *.go) $$(call rwildcard, ./core, *.go) $$(call rwildcard, ./listener, *.go) $$(call rwildcard, ./writer, *.go)
	$(CC) $(DFLAGS) $(CFLAGS) -o $(BUILD_DIR)/fossil fossil.go

.PHONY: release
release: fossil.go $$(call rwildcard, ./cmd, *.go) $$(call rwildcard, ./core, *.go) $$(call rwildcard, ./listener, *.go) $$(call rwildcard, ./writer, *.go)
	$(CC) $(CFLAGS) -ldflags "-s -w" -o $(BUILD_DIR)/fossil fossil.go

.PHONY: dist
dist: fossil.go $$(call rwildcard, ./cmd, *.go) $$(call rwildcard, ./core, *.go) $$(call rwildcard, ./listener, *.go) $$(call rwildcard, ./writer, *.go)
	$(CROSS) $(CC) $(CFLAGS) -ldflags "-s -w" -o $(BUILD_DIR)/fossil fossil.go

.PHONY: install
install: 
	install -m 0755 $(BUILD_DIR)/fossil $(PREFIX)/bin

.PHONY: uninstall
uninstall: 
	rm -f $(PREFIX)/bin/fossil

.PHONY: lint
lint:
	@command -v gometalinter >/dev/null 2>&1 || { echo >&2 "gometalinter is required but not available please follow instructions from https://github.com/alecthomas/gometalinter"; exit 1; }
	gometalinter --deadline=180s --disable-all --enable=gofmt ./listener ./writer ./cmd/... ./core/... ./
	gometalinter --deadline=180s --disable-all --enable=vet ./listener ./writer ./cmd/... ./core/... ./
	gometalinter --deadline=180s --disable-all --enable=golint ./listener ./writer ./cmd/... ./core/... ./
	gometalinter --deadline=180s --disable-all --enable=ineffassign ./listener ./writer ./cmd/... ./core/... ./
	gometalinter --deadline=180s --disable-all --enable=misspell ./listener ./writer ./cmd/... ./core/... ./
	gometalinter --deadline=180s --disable-all --enable=staticcheck ./listener ./writer ./cmd/... ./core/... ./

.PHONY: format
format:
	gofmt -w -s ./cmd ./core ./listener ./writer fossil.go

.PHONY: dev
dev: format lint build

.PHONY: clean
clean:
	-rm -r build
