FROM golang:1.26.2 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download 

COPY . . 
RUN mkdir -p /out && go build -o /out/api ./cmd

FROM alpine:3.24 as runtime

// Alpine does not include TLS certificates to make HTTPS requests
RUN apk add --no-cache ca-certificates

COPY --from=build /out/api /app/api


CMD ["app/api"]