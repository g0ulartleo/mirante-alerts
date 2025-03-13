FROM golang:1.23-alpine

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev make

ENV CGO_ENABLED=1

COPY . .

RUN go mod download

RUN make go-install-templ
RUN make go-install-air

RUN make build

EXPOSE 40169
