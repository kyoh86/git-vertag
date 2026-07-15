# Usage:
# docker run -it --rm ghcr.io/kyoh86/git-vertag:latest -v ${PWD}:/work patch
FROM golang:alpine AS build-stage
ADD . /work
WORKDIR /work
RUN go build -o git-vertag .

FROM alpine:latest
LABEL org.opencontainers.image.source=https://github.com/kyoh86/git-vertag
LABEL org.opencontainers.image.description="A tool to manage version-tag with the semantic versioning specification."
LABEL org.opencontainers.image.licenses=MIT
RUN apk update && apk add git
WORKDIR /work
COPY --from=build-stage /work/git-vertag /usr/local/bin/git-vertag
ENTRYPOINT ["/usr/local/bin/git-vertag"]
