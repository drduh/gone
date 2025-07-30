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

SERVICE   = $(APPNAME).service

ASSET_CSS = assets/style.css
SETTINGS  = settings/defaultSettings.json

CONF_DIR ?= /etc/$(APPNAME)
DEST_BIN  = /usr/local/bin/$(APPNAME)
DEST_CONF = $(CONF_DIR)/config
DEST_CSS  = $(CONF_DIR)/$(ASSET_CSS)
DEST_SERV = /etc/systemd/system/$(SERVICE)

MOD_BIN   = 0755
MOD_FILE  = 0644

TESTCOVER = testCoverage

all: fmt test build

prep:
	@mkdir -p $(OUT)

build-linux:
	@$(BLDLINUX)

build: prep build-linux

run: build
	@$(OUT)/$(BINLINUX)

debug: build
	@$(OUT)/$(BINLINUX) -debug

version: build
	@$(OUT)/$(BINLINUX) -version

install: install-assets install-config install-bin install-service reload-service

install-assets:
	@sudo install -Dm $(MOD_FILE) $(ASSET_CSS) $(DEST_CSS)
	@printf "Installed $(DEST_CSS)\n"

install-config:
	@sudo install -Dm $(MOD_FILE) $(SETTINGS) $(DEST_CONF)
	@printf "Installed $(DEST_CONF)\n"

install-bin: build
	@sudo install -Dm $(MOD_BIN) $(OUT)/$(BINLINUX) $(DEST_BIN)
	@printf "Installed $(DEST_BIN)\n"

install-service:
	@sudo install -Dm $(MOD_FILE) $(SERVICE) $(DEST_SERV)
	@sudo systemctl enable $(SERVICE)
	@printf "Installed $(DEST_SERV)\n"

reload-service:
	@printf "Restarting services ...\n"
	@sudo systemctl daemon-reload
	@sudo systemctl restart $(SERVICE)

fmt:
	@$(GO) fmt ./...

test:
	@$(GO) test ./...

test-verbose:
	@$(GO) test -v ./...

test-cover:
	@$(GO) test -coverprofile=$(TESTCOVER) ./...

cover: test-cover
	@$(GO) tool cover -html=$(TESTCOVER) -o $(TESTCOVER).html

doc:
	@$(GODOC) -http :8000

clean:
	@rm -rf $(OUT) $(TESTCOVER) $(TESTCOVER).html
