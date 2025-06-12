ROOT      = gone
APPNAME  ?= $(ROOT)
AUTHOR   ?= drduh
GIT      ?= github.com/$(AUTHOR)
VERSION  ?= $(shell date +"%Y.%m.%d")

CMD       = cmd
SRC       = $(CMD)/main.go
OUT       = release

GO       ?= go
GODOC    ?= ${HOME}/go/bin/godoc

BUILDPKG  = $(GIT)/$(APPNAME)/version
BUILDARCH = $(shell $(GO) env GOHOSTARCH)
BUILDVERS = $(shell $(GO) env GOVERSION)
BUILDOS   = $(shell $(GO) env GOHOSTOS)
BUILDTIME = $(shell date +"%Y-%m-%dT%H:%M:%S")
BUILDFLAG = \
  -X "$(BUILDPKG).Arch=$(BUILDARCH)" \
  -X "$(BUILDPKG).Go=$(BUILDVERS)" \
  -X "$(BUILDPKG).Host=$(shell hostname -f)" \
  -X "$(BUILDPKG).Id=$(APPNAME)" \
  -X "$(BUILDPKG).Path=$(shell pwd)" \
  -X "$(BUILDPKG).OS=$(BUILDOS)" \
  -X "$(BUILDPKG).Time=$(BUILDTIME)" \
  -X "$(BUILDPKG).User=$(shell whoami)" \
  -X "$(BUILDPKG).Version=$(VERSION)"
BUILDCMD  = $(GO) build -ldflags '-s -w $(BUILDFLAG)'
BINLINUX  = $(APPNAME)-$(BUILDOS)-$(BUILDARCH)-$(VERSION)
BLDLINUX  = GOOS=$(BUILDOS) GOARCH=$(BUILDARCH) \
            $(BUILDCMD) -o $(OUT)/$(BINLINUX) $(SRC)

INSTALL  ?= /usr/local/bin/$(APPNAME)

TESTCOVER = testCoverage

all: fmt test build

run: build
	@$(OUT)/$(BINLINUX)

debug: build
	@$(OUT)/$(BINLINUX) -debug

version: build
	@$(OUT)/$(BINLINUX) -version

build: prep linux

linux:
	@$(BLDLINUX)

prep:
	@mkdir -p $(OUT)

install: build
	@sudo install -Dm755 $(OUT)/$(BINLINUX) $(INSTALL)
	@printf "Installed $(INSTALL)\n"

fmt:
	@$(GO) fmt ./...

test:
	@$(GO) test ./...

test-verbose:
	@$(GO) test -v ./...

test-cover:
	@$(GO) test -coverprofile=$(TESTCOVER) ./...

cover: test-cover
	@$(GO) tool cover -html=$(TESTCOVER) \
		-o $(TESTCOVER).html

doc:
	@$(GODOC) -http :8000

clean:
	@rm -rf $(OUT) \
		$(TESTCOVER) $(TESTCOVER).html
