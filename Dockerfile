FROM golang:1.17-alpine as golang

FROM alpine:latest

COPY --from=golang /usr/local/go/lib/time/zoneinfo.zip /
COPY --from=golang /app/main /app/main

ENV ZONEINFO=/zoneinfo.zip

WORKDIR  /app

ENTRYPOINT ./main