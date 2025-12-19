FROM golang:1.25-alpine3.23 AS builder

ENV CGO_ENABLED=0

RUN apk update && apk add --no-cache git tzdata ca-certificates
RUN mkdir /app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o /go/bin/app \
   "./cmd/server"

FROM alpine:latest

WORKDIR /app

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/bin/app ./app_binary
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/docs ./docs

EXPOSE 9898

CMD ["./app_binary"]