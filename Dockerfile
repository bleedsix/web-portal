FROM golang:1.20-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/bleedsix/web-portal
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/web-portal /go/src/github.com/bleedsix/web-portal


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/web-portal /usr/local/bin/web-portal
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["web-portal"]
