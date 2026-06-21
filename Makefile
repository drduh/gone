# https://github.com/drduh/gone/blob/main/Makefile
ROOT      = gone

APPNAME  ?= $(ROOT)
APPVERS  ?= $(shell date +"%Y.%m.%d")
AUTHOR   ?= drduh
GITNAME  ?= github.com
GITREPO  ?= $(GITNAME)/$(AUTHOR)

PKG       = ./...
CMD       = cmd
SRC       = $(CMD)/main.go
OUT       = release

GOCMD    ?= go
GODOC    ?= ${HOME}/go/bin/godoc
GOLINT   ?= golangci-lint
GOLINTARG =
GOSEC    ?= gosec
GOSTATIC ?= staticcheck

CONTAIN  ?= container
DOCKER   ?= docker

BUILDARCH = $(shell $(GOCMD) env GOHOSTARCH)
BUILDGIT  = $(shell git log -1 --format=%H \
            2>/dev/null || printf "unknown")
BUILDHOST = $(shell hostname -f)
BUILDOS   = $(shell $(GOCMD) env GOHOSTOS)
BUILDPATH = $(shell pwd)
BUILDTIME = $(shell date +"%Y-%m-%dT%H:%M:%S")
BUILDUSER = $(shell whoami)
BUILDVERS = $(shell $(GOCMD) env GOVERSION)
VERSPKG   = $(GITREPO)/$(APPNAME)/version
BUILDFLAG = -X "$(VERSPKG).Arch=$(BUILDARCH)" \
            -X "$(VERSPKG).Commit=$(BUILDGIT)" \
            -X "$(VERSPKG).Go=$(BUILDVERS)" \
            -X "$(VERSPKG).Host=$(BUILDHOST)" \
            -X "$(VERSPKG).ID=$(APPNAME)" \
            -X "$(VERSPKG).Path=$(BUILDPATH)" \
            -X "$(VERSPKG).System=$(BUILDOS)" \
            -X "$(VERSPKG).Time=$(BUILDTIME)" \
            -X "$(VERSPKG).User=$(BUILDUSER)" \
            -X "$(VERSPKG).Version=$(APPVERS)"

# example - gone-darwin-arm64-2026.12.31
BINNAME  ?= $(APPNAME)-$(BUILDOS)-$(BUILDARCH)-$(APPVERS)
CMDBUILD  = $(GOCMD) build -trimpath \
            -ldflags '-s -w $(BUILDFLAG)'
GOBUILD   = GOOS=$(BUILDOS) GOARCH=$(BUILDARCH) \
            $(CMDBUILD) \
            -o "$(OUT)/$(BINNAME)" "$(SRC)"
GORACE    = GOOS=$(BUILDOS) GOARCH=$(BUILDARCH) \
            $(CMDBUILD) -race \
            -o "$(OUT)/$(BINNAME)-race" "$(SRC)"

SERVICE   = $(APPNAME).service
SYSTEMCTL = systemctl

ASSETS    = assets
ASSET_CSS = $(ASSETS)/style.css
SETTINGS  = settings/defaultSettings.json

CONF_DIR ?= /etc/$(APPNAME)
DEST_BIN  = /usr/local/bin/$(APPNAME)
DEST_CONF = $(CONF_DIR)/config
DEST_CSS  = $(CONF_DIR)/$(ASSET_CSS)
DEST_SERV = /etc/systemd/system/$(SERVICE)

MOD_EXEC  = 0755
MOD_FILE  = 0644

TESTARGS ?=
TESTCOVER = testCoverage
TIMEOUT  ?= 1m
CMDTEST   = $(GOCMD) test -trimpath
CMDCOVER  = $(CMDTEST) \
            -coverprofile=$(TESTCOVER) $(PKG)

WARN      = tput setaf 3 ; printf "%s\n" "${1}" ; \
            tput sgr0

all: fmt lint test build

prep-build:
	@mkdir -p $(OUT)

build: prep-build
	@$(GOBUILD)

run: build
	@$(OUT)/$(BINNAME)

run-container: build-container
	@$(CONTAIN) run $(APPNAME)-$(APPVERS)

debug: build
	@$(OUT)/$(BINNAME) -debug

version: build
	@$(OUT)/$(BINNAME) -version

release: build
	@printf "built release: %s\n" \
		"$$(file "$(OUT)/$(BINNAME)")"

prep-container:
	@$(CONTAIN) system start

build-container: prep-container
	@$(CONTAIN) build -t $(APPNAME)-$(APPVERS) .

install: install-assets install-bin \
	install-logdir \
	install-service reload-service \
	check-service

install-assets:
	@printf "Installing $(DEST_CSS) ... "
	@sudo install -Dm $(MOD_FILE) \
		$(ASSET_CSS) $(DEST_CSS)
	@printf "done\n"

install-bin: build
	@printf "Installing $(BINNAME) to $(DEST_BIN) ... "
	@sudo install -Dm $(MOD_EXEC) \
		-o root -g $(APPNAME) $(OUT)/$(BINNAME) $(DEST_BIN)
	@printf "done\n"

install-logdir:
	@printf "Installing /var/log/$(APPNAME) ... "
	@sudo install -Dm $(MOD_EXEC) \
		-o $(APPNAME) -g $(APPNAME) -d /var/log/$(APPNAME)
	@printf "done\n"

install-user:
	@id -u $(APPNAME) > /dev/null 2>&1 || \
		sudo useradd --system --no-create-home \
		--shell /usr/sbin/nologin $(APPNAME)

install-config: install-user
	@printf "Installing $(DEST_CONF) ... "
	@sudo install -Dm $(MOD_FILE) \
		-o root -g $(APPNAME) $(SETTINGS) $(DEST_CONF)
	@printf "done\n"

install-service: install-config
	@printf "Installing $(DEST_SERV) ... "
	@sudo install -Dm $(MOD_FILE) \
		$(SERVICE) $(DEST_SERV)
	@sudo $(SYSTEMCTL) enable $(SERVICE)
	@printf "done\n"

reload-service:
	@printf "Restarting services ... "
	@sudo $(SYSTEMCTL) daemon-reload
	@sudo $(SYSTEMCTL) restart $(SERVICE)
	@printf "done\n"

check-service:
	@printf "Checking service install ... \n"
	@sleep 2
	@$(SYSTEMCTL) status $(APPNAME) || \
		$(DEST_BIN) -conf $(DEST_CONF)

uninstall:
	@sudo $(SYSTEMCTL) stop $(APPNAME)
	@sudo $(SYSTEMCTL) disable $(APPNAME)
	@sudo rm -f $(DEST_SERV)

fmt:
	@$(GOCMD) fmt $(PKG)

test-race:    TESTARGS = -race
test-short:   TESTARGS = -short
test-verbose: TESTARGS = -v

test test-race test-short test-verbose:
	@$(CMDTEST) $(TESTARGS) -timeout=$(TIMEOUT) $(PKG)

RUN_IF_FOUND = if command -v $(1) >/dev/null 2>&1 ; \
		then $(1) $(2) ; else \
		$(call WARN,skipping '$@': '$(1)' not found); fi

lint-verbose: GOLINTARG = --verbose

lint lint-verbose:
	@printf "linting ... "
	@$(call RUN_IF_FOUND,$(GOLINT),run $(GOLINTARG) $(PKG))

sec:
	@$(call RUN_IF_FOUND,$(GOSEC),$(PKG))

static:
	@$(call RUN_IF_FOUND,$(GOSTATIC),$(PKG))

build-race: prep-build
	@$(GORACE)

race: build-race
	@$(OUT)/$(BINNAME)-race -debug

clean: clean-coverage
	@rm -rf $(OUT)

clean-coverage:
	@rm -rf $(TESTCOVER) $(TESTCOVER).html

clean-cache:
	@$(GOCMD) clean -cache -testcache -modcache
	@$(GOLINT) cache clean

cover: test-cover
	@$(GOCMD) tool cover \
		-html="$(TESTCOVER)" -o "$(TESTCOVER).html"
	@printf "total test coverage: %s" \
		"$$($(GOCMD) tool cover -func="$(TESTCOVER)" | \
		awk '/^total:/{print $$3}')"
	@printf " - see %s\n" "$(TESTCOVER).html"

test-cover:
	@$(CMDCOVER)

doc:
	@$(GODOC) -http :8000

c: clean
celan: clean
clena: clean
coen: coverage
coveage: coverage
coverae: coverage
coverage: cover
d: debug
devug: debug
f: fmt
format: fmt
gosec: sec
litn: lint
prep: prep-build
prod: release
r: run
restart: reload-service
restart-service: reload-service
t: test
tset: test
un: run
urn: run
v: verbose
verbose: debug
