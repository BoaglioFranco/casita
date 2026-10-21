FROM golang:1.26.2 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download 

COPY . . 
RUN mkdir -p /out && go build -o /out/api ./cmd

FROM debian:bookworm-slim as runtime

COPY --from=build /out/api /app/api


CMD ["/app/api"]