ROOT      = gone
AUTHOR   ?= drduh

APP      ?= $(ROOT)
GIT      ?= github.com/$(AUTHOR)
GO       ?= go
VERSION  ?= $(shell date +"%Y.%m.%d")

CMD       = cmd
SRC       = $(CMD)/main.go
OUT       = release

BINARY    = $(APP)-$(VERSION)
BUILDUSER = $(shell whoami)
BUILDOS   = $(shell $(GO) env GOHOSTOS)
BUILDARCH = $(shell $(GO) env GOHOSTARCH)
BUILDVERS = $(shell $(GO) env GOVERSION)
BUILDPKG  = $(GIT)/$(APP)/version
BUILDFLAG = -X "$(BUILDPKG).Id=$(APP)" \
            -X "$(BUILDPKG).Version=$(VERSION)" \
            -X "$(BUILDPKG).User=$(BUILDUSER)" \
            -X "$(BUILDPKG).OS=$(BUILDOS)" \
            -X "$(BUILDPKG).Arch=$(BUILDARCH)" \
            -X "$(BUILDPKG).Go=$(BUILDVERS)"
BUILDCMD  = $(GO) build -ldflags '-s -w $(BUILDFLAG)'
BINLINUX  = $(BINARY)-$(BUILDOS)-$(BUILDARCH)
BLDLINUX  = GOOS=$(BUILDOS) GOARCH=$(BUILDARCH) \
            $(BUILDCMD) -o $(OUT)/$(BINLINUX) $(SRC)

all: fmt dev

dev: build
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

fmt:
	@$(GO) fmt ./...

test:
	@$(GO) test ./...

clean:
	@rm -rf $(OUT)
