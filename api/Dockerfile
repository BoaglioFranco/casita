FROM golang:1.26.2 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download 

COPY . . 
RUN mkdir -p /out && go build -o /out/api ./cmd

FROM debian:bookworm-slim as runtime

RUN apt-get update \
    && apt-get install -y --no-install-recommends sqlite3 \
    && rm -rf /var/lib/apt/lists/*
    
COPY migrations /app/migrations


COPY --from=build /out/api /app/api


CMD ["/app/api"]