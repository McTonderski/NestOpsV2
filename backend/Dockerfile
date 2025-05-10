### First stage: build with CGO + SQLite support
FROM golang:1.24.1-alpine3.20 AS builder

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=arm64

WORKDIR /app

# Install musl & sqlite headers
RUN apk add --no-cache build-base sqlite-dev

COPY . .

RUN go mod download

# 👇 Force fully static build (no dynamic libs)
RUN go build -tags 'osusergo netgo static_build' -ldflags '-linkmode external -extldflags "-static"' -o app .

### Final tiny image
FROM scratch

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/static ./static
COPY --from=builder /app/.env .env
COPY --from=builder /app/services ./services

ENTRYPOINT ["./app"]