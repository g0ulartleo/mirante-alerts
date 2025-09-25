FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache build-base sqlite-dev make git

ENV CGO_ENABLED=1

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN make build

FROM golang:1.23-alpine AS dev

WORKDIR /app

RUN apk add --no-cache build-base sqlite-dev git

ENV CGO_ENABLED=1

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest && \
    go install github.com/air-verse/air@v1.61.7

EXPOSE 40169

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata sqlite-libs && update-ca-certificates

COPY --from=builder /app/bin/ /app/bin/
COPY --from=builder /app/static /app/static
COPY --from=builder /app/config /app/config

ENV HTTP_ADDR=0.0.0.0

EXPOSE 40169

CMD ["./bin/http-server"]
