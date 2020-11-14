FROM golang:alpine AS build-stage
ADD . /work
WORKDIR /work
RUN go build -o git-vertag .

FROM alpine:latest
COPY --from=build-stage /work/git-vertag /usr/local/bin/git-vertag
ENTRYPOINT ["/usr/local/bin/git-vertag"]
