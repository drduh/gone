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
GOSEC    ?= gosec
GOSTATIC ?= staticcheck

CONTAIN  ?= container
DOCKER   ?= docker

BUILDARCH = $(shell $(GOCMD) env GOHOSTARCH)
BUILDGIT  = $(shell git log -1 --format=%h \
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
            -X "$(VERSPKG).Id=$(APPNAME)" \
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

TESTCOVER = testCoverage
CMDTEST   = $(GOCMD) test -trimpath
CMDCOVER  = $(CMDTEST) \
            -coverprofile=$(TESTCOVER) $(PKG)

TIMEOUT  ?= 1m

WARN      = tput setaf 3 ; \
            printf "%s\n" "${1}" ; \
            tput sgr0

all: fmt build test lint

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
		"$$(file $(OUT)/$(BINNAME))"

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

uninstall-service:
	@sudo $(SYSTEMCTL) stop $(APPNAME)
	@sudo $(SYSTEMCTL) disable $(APPNAME)
	@sudo rm -f /etc/systemd/system/$(APPNAME)

fmt:
	@$(GOCMD) fmt $(PKG)

test:
	@$(CMDTEST) -timeout=$(TIMEOUT) $(PKG)

test-race:
	@$(CMDTEST) -race -timeout=$(TIMEOUT) $(PKG)

test-short:
	@$(CMDTEST) -short -timeout=$(TIMEOUT) $(PKG)

test-verbose:
	@$(CMDTEST) -v -timeout=$(TIMEOUT) $(PKG)

test-cover:
	@$(CMDCOVER)

test-cover-total: test-cover
	@echo "total coverage: \
		$$($(GOCMD) tool cover -func=$(TESTCOVER) | \
		grep total: | awk '{print $$3}')"

test-cover-all: test-cover-total

lint:
	@if command -v $(GOLINT) >/dev/null 2>&1 ; then \
		$(GOLINT) run $(PKG) ; else \
		$(call WARN,skipping '$@': '$(GOLINT)' not found); \
	fi

lint-verbose:
	@if command -v $(GOLINT) >/dev/null 2>&1 ; then \
		$(GOLINT) run --verbose $(PKG) ; else \
		$(call WARN,skipping '$@': '$(GOLINT)' not found); \
	fi

sec:
	@if command -v $(GOSEC) >/dev/null 2>&1 ; then \
		$(GOSEC) run $(PKG) ; else \
		$(call WARN,skipping '$@': '$(GOSEC)' not found); \
	fi

static:
	@if command -v $(GOSTATIC) >/dev/null 2>&1 ; then \
		$(GOSTATIC) $(PKG) ; else \
		$(call WARN,skipping '$@': '$(GOSTATIC)' not found); \
	fi

build-race: prep
	@$(GORACE)

race: build-race
	@$(OUT)/$(BINNAME)-race -debug

cover: test-cover
	@$(GOCMD) tool cover \
		-html=$(TESTCOVER) -o $(TESTCOVER).html
	@printf "cover: %s\n" \
		"$$(file $(TESTCOVER).html)"

doc:
	@$(GODOC) -http :8000

clean: clean-coverage
	@rm -rf $(OUT)

clean-coverage:
	@rm -rf $(TESTCOVER) $(TESTCOVER).html

clean-cache:
	@$(GOCMD) clean -cache -testcache -modcache
	@$(GOLINT) cache clean

c: clean
celan: clean
clena: clean
coveage: coverage
coverae: coverage
coverage: cover
d: debug
devug: debug
f: fmt
prod: release
t: test
tset: test
r: run
urn: run
v: verbose
verbose: debug
