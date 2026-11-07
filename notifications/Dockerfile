FROM golang:1.26.2 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir -p /out \
    && go build -o /out/notifications ./cmd \
    && go build -o /out/migrate ./cmd/migrate

FROM debian:bookworm-slim AS runtime
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /out/notifications /app/notifications
COPY --from=build /out/migrate /app/migrate

CMD ["/app/notifications"]
