# https://github.com/drduh/gone/blob/main/Makefile
ROOT      = gone

APPNAME  ?= $(ROOT)
APPGEN   ?= v1
APPVERS  ?= $(APPGEN).$(shell date +"%Y.%m.%d")
AUTHOR   ?= drduh
GITNAME  ?= github.com
GITREPO  ?= $(GITNAME)/$(AUTHOR)

ARG       =
OUT       = release
PKG       = ./...
SRC       = cmd/main.go

CMD_GO   ?= go
GODOC    ?= ${HOME}/go/bin/godoc
GOLINT   ?= golangci-lint
GOLINTARG =
GOSEC    ?= gosec
GOSTATIC ?= staticcheck

CONTAIN  ?= container
DOCKER   ?= docker

BUILDARCH = $(shell $(CMD_GO) env GOHOSTARCH)
BUILDGIT  = $(shell git log -1 --format=%H \
            2>/dev/null || printf "unknown")
BUILDHOST = $(shell hostname -f)
BUILDOS   = $(shell $(CMD_GO) env GOHOSTOS)
BUILDPATH = $(shell pwd)
BUILDTIME = $(shell date +"%Y-%m-%dT%H:%M:%S")
BUILDUSER = $(shell whoami)
BUILDVERS = $(shell $(CMD_GO) env GOVERSION)
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

# example - gone-darwin-arm64-v1.2026.12.31
BINNAME  ?= $(APPNAME)-$(BUILDOS)-$(BUILDARCH)-$(APPVERS)
BINRACE   = $(BINNAME)-race
BUILDCMD  = $(CMD_GO) build -trimpath -ldflags '-s -w $(BUILDFLAG)'

BUILDBASE = GOOS=$(BUILDOS) GOARCH=$(BUILDARCH) $(BUILDCMD)
CMD_BUILD = $(BUILDBASE) -o "$(OUT)/$(BINNAME)" "$(SRC)"
CMD_RACE  = $(BUILDBASE) -race -o "$(OUT)/$(BINRACE)" "$(SRC)"

SERVICE   = $(APPNAME).service
SYSTEMCTL = systemctl

ASSET_DIR = assets
ASSET_CSS = $(ASSET_DIR)/style.css
SETTINGS  = settings/defaultSettings.json

CONF_DIR ?= /etc/$(APPNAME)
DEST_BIN  = /usr/local/bin/$(APPNAME)
DEST_CONF = $(CONF_DIR)/config
DEST_CSS  = $(CONF_DIR)/$(ASSET_CSS)
DEST_SERV = /etc/systemd/system/$(SERVICE)

MOD_EXEC  = 0755
MOD_FILE  = 0644

ARG_TEST ?=
TESTCOVER = testCoverage
TESTTIME ?= 1m

CMD_TEST  = $(CMD_GO) test -trimpath
CMD_COVER = $(CMD_TEST) \
            -coverprofile=$(TESTCOVER) $(PKG)

AUTHCRED ?= mySecret

WARN      = tput setaf 3 ; printf "%s\n" "${1}" ; tput sgr0

all: fmt lint test build

prep-build:
	@mkdir -p $(OUT)

build: prep-build
	@$(CMD_BUILD)

debug:   ARG += -debug
version: ARG += -version

run debug version: build
	@$(OUT)/$(BINNAME) -auth $(AUTHCRED) $(ARG)

release: build
	@printf "built release: %s\n" \
		"$$(file "$(OUT)/$(BINNAME)")"

prep-container:
	@$(CONTAIN) system start

build-container: prep-container
	@$(CONTAIN) build -t $(APPNAME)-$(APPVERS) .

run-container: build-container
	@$(CONTAIN) run $(APPNAME)-$(APPVERS)

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
		$(DEST_BIN) -config $(DEST_CONF)

uninstall:
	@sudo $(SYSTEMCTL) stop $(APPNAME)
	@sudo $(SYSTEMCTL) disable $(APPNAME)
	@sudo rm -f $(DEST_SERV)

fmt:
	@$(CMD_GO) fmt $(PKG)

test-race:    ARG_TEST = -race
test-short:   ARG_TEST = -short
test-verbose: ARG_TEST = -v

test test-race test-short test-verbose:
	@$(CMD_TEST) $(ARG_TEST) -timeout=$(TESTTIME) $(PKG)

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
	@$(CMD_RACE)

race: build-race
	@$(OUT)/$(BINRACE) -debug

clean: clean-cert clean-coverage
	@rm -rf $(OUT)

clean-cert:
	@rm -rf cert.pem key.pem

clean-coverage:
	@rm -rf $(TESTCOVER) $(TESTCOVER).html

clean-cache:
	@$(CMD_GO) clean -cache -testcache -modcache
	@$(GOLINT) cache clean

cover: test-cover
	@$(CMD_GO) tool cover \
		-html="$(TESTCOVER)" -o "$(TESTCOVER).html"
	@printf "total test coverage: %s" \
		"$$($(CMD_GO) tool cover -func="$(TESTCOVER)" | \
		awk '/^total:/{print $$3}')"
	@printf " - see %s\n" "$(TESTCOVER).html"

test-cover:
	@$(CMD_COVER)

view-cover: cover
	@open "$(TESTCOVER).html"

doc:
	@$(GODOC) -http :8000

key.pem:
	@openssl ecparam -name prime256v1 -genkey -noout -out $@

cert.pem: key.pem
	@openssl req -new -x509 \
		-key $< -out $@ \
		-days 8 -subj "/CN=${APPNAME}.${APPVERS}"

b: build
c: clean
cealn: clean
celan: clean
cert: cert.pem
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
open-cover: view-cover
pem: cert
prep: prep-build
prod: release
r: run
restart: reload-service
restart-service: reload-service
rin: run
t: test
tets: test
tset: test
un: run
urn: run
v: verbose
vers: version
verbose: debug
verison: version
