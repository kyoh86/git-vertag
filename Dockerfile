# Usage:
# docker run -it --rm gcr.io/kyoh86/git-vertag:latest -v ${PWD}:/work patch
FROM golang:alpine AS build-stage
ADD . /work
WORKDIR /work
RUN go build -o git-vertag .

FROM alpine:latest
RUN apk update && apk add git
WORKDIR /work
COPY --from=build-stage /work/git-vertag /usr/local/bin/git-vertag
ENTRYPOINT ["/usr/local/bin/git-vertag"]
