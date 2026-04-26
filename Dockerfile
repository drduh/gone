FROM golang:1.26 AS build
WORKDIR /src
COPY . .
ARG APPNAME=gone
ARG AUTHOR=drduh
RUN set -eu; \
  BUILDPKG="github.com/${AUTHOR}/${APPNAME}/version"; \
  BUILDARCH="$(go env GOHOSTARCH)"; \
  BUILDVERS="$(go env GOVERSION)"; \
  BUILDOS="$(go env GOHOSTOS)"; \
  BUILDGIT="$(git log -1 --format=%h 2>/dev/null || printf '0000000')"; \
  BUILDTIME="$(date -u +'%Y-%m-%dT%H:%M:%S')"; \
  VERSION="${VERSION:-$(date +'%Y.%m.%d')}"; \
  HOST="${HOST:-$(hostname 2>/dev/null || true)}"; \
  BUILDUSER="${BUILDUSER:-$(id -un 2>/dev/null || true)}"; \
  BUILDFLAG="\
    -X ${BUILDPKG}.Arch=${BUILDARCH} \
    -X ${BUILDPKG}.Commit=${BUILDGIT} \
    -X ${BUILDPKG}.Go=${BUILDVERS} \
    -X ${BUILDPKG}.Host=${HOST} \
    -X ${BUILDPKG}.Id=${APPNAME} \
    -X ${BUILDPKG}.Path=$(pwd) \
    -X ${BUILDPKG}.System=${BUILDOS} \
    -X ${BUILDPKG}.Time=${BUILDTIME} \
    -X ${BUILDPKG}.User=${BUILDUSER} \
    -X ${BUILDPKG}.Version=${VERSION}"; \
  go build -trimpath -ldflags "-s -w ${BUILDFLAG}" -o /out/gone ./cmd

FROM gcr.io/distroless/base-debian13:nonroot
WORKDIR /etc/gone
COPY --from=build /out/gone /usr/local/bin/gone
COPY --from=build /src/assets/ ./assets/
ENTRYPOINT ["/usr/local/bin/gone"]
