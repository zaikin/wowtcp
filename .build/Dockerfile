FROM golang:1.23 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags '-w -s' -a -o /build/app ./cmd/wowtcp

FROM alpine:3.14

RUN apk update
RUN apk upgrade

WORKDIR /app

COPY --from=builder /build/app /app/app
ENTRYPOINT [ "/app/app" ]
