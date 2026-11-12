# Build stage: compile a static Gateway binary.
FROM golang:1.25.14-alpine3.24 AS build
WORKDIR /src

# Download modules first so this layer is cached until go.mod/go.sum change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/gateway ./cmd/server

# Run stage: minimal image with the binary only.
FROM alpine:3.24.4
COPY --from=build /out/gateway /usr/local/bin/gateway

# The image contains no route configuration (D-17): mount the config file
# at /etc/gateway/config.yaml (read-only). The server listens on :8080.
USER nobody
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/gateway"]
CMD ["-config", "/etc/gateway/config.yaml"]
