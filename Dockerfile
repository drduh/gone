FROM golang:1.26 AS build
WORKDIR /src
COPY . .
RUN chmod -R 755 assets
ARG APPNAME=gone
ARG AUTHOR=drduh
ARG GITNAME=github.com
RUN set -eu; \
  BUILDARCH="$(go env GOHOSTARCH)"; \
  BUILDGIT="$(git log -1 --format=%H 2>/dev/null || printf '0000000')"; \
  BUILDHOST="$(hostname 2>/dev/null || true)"; \
  BUILDOS="$(go env GOHOSTOS)"; \
  BUILDPATH="$(pwd)"; \
  VERSPKG="${GITNAME}/${AUTHOR}/${APPNAME}/version"; \
  BUILDTIME="$(date -u +'%Y-%m-%dT%H:%M:%S')"; \
  BUILDUSER="$(id -un 2>/dev/null || true)"; \
  BUILDVERS="$(go env GOVERSION)"; \
  VERSION="$(date +'%Y.%m.%d')"; \
  BUILDFLAG="-X ${VERSPKG}.Arch=${BUILDARCH} \
             -X ${VERSPKG}.Commit=${BUILDGIT} \
             -X ${VERSPKG}.Go=${BUILDVERS} \
             -X ${VERSPKG}.Host=${BUILDHOST} \
             -X ${VERSPKG}.ID=${APPNAME} \
             -X ${VERSPKG}.Path=${BUILDPATH} \
             -X ${VERSPKG}.System=${BUILDOS} \
             -X ${VERSPKG}.Time=${BUILDTIME} \
             -X ${VERSPKG}.User=${BUILDUSER} \
             -X ${VERSPKG}.Version=${VERSION}"; \
  go build -trimpath -ldflags "-s -w ${BUILDFLAG}" -o /out/gone ./cmd

FROM gcr.io/distroless/base-debian13:nonroot
WORKDIR /etc/gone
COPY --from=build /out/gone /usr/local/bin/gone
COPY --from=build /src/assets/ ./assets/
#ENTRYPOINT ["/usr/local/bin/gone", "-debug"]
#ENTRYPOINT ["/usr/local/bin/gone", "-version"]
ENTRYPOINT ["/usr/local/bin/gone"]
