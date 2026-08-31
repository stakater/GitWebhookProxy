# A SELF-CONTAINED MULTI-STAGE BUILD, WHICH THE UPSTREAM IMAGES WERE NOT.
#
# `build/package/Dockerfile` is three lines that `COPY ./gitwebhookproxy /` — it
# expects a binary already compiled beside it by the Makefile, so it cannot be
# handed to any platform that builds from a repository. `Dockerfile.build` is a
# Go 1.13 image driving `packr` and a GOPATH layout that modules replaced years
# ago.
#
# This builds from source in one pass. `docker build .` produces a runnable
# image, and so does any platform that points at this repository — which is the
# whole reason it exists.

# PINNED BY DIGEST, NEVER `:alpine`. A moving tag lets the toolchain that
# compiles a network-facing binary change under a running deployment on somebody
# else's release schedule, and "it rebuilt itself differently" is not a thing to
# debug at two in the morning.
FROM golang:1.25-alpine@sha256:1ae0735f00daffa3aaf1363a5184c0d2dc55c78e3db4ec70241cdac97bf84b59 AS build

WORKDIR /src

# Dependencies first, so a source-only change does not re-resolve the module
# graph. `go mod download` is the whole of it — the build is stdlib plus three
# small packages.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO OFF AND A STATIC BINARY, so the runtime stage can be `scratch`-adjacent
# and carry no libc of its own to keep patched. `-trimpath` keeps build-machine
# paths out of the binary; `-s -w` drop the symbol and DWARF tables, which are
# of no use in production and are pure attack surface in a network daemon.
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -trimpath -ldflags="-s -w" -o /out/gitwebhookproxy .

# THE RUNTIME STAGE CARRIES NO SHELL AND NO PACKAGE MANAGER.
#
# This process is reachable from the internet by design — that is its entire
# job. `static-debian12` gives it CA certificates and nothing else: no `sh`, no
# `apk`, no `curl`, so a command-injection bug has nothing to inject into.
# Upstream used `stakater/base-alpine:3.7`, an Alpine 3.7 base last updated in
# 2018 that ships a full shell and package manager.
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab

# NON-ROOT, and the image's own user rather than one created here — there is no
# shell to run `useradd` with, which is the point.
USER nonroot:nonroot

COPY --from=build /out/gitwebhookproxy /gitwebhookproxy

# Documentation rather than a binding; the listen address comes from
# `GWP_LISTENADDRESS`.
EXPOSE 8080

ENTRYPOINT ["/gitwebhookproxy"]
