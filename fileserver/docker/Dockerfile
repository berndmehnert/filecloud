# Build stage
FROM golang:1.24.4 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/fileserver ./cmd/fileserver

# Runtime stage
FROM scratch
# Optional: add CA certs for HTTPS clients; use distroless or alpine if you need more tools.
# For minimal image with certs, uncomment lines below and use gcr.io/distroless/...
# FROM gcr.io/distroless/cc-debian11

# copy binary from build
COPY --from=build /out/fileserver /usr/local/bin/fileserver

# Create writable storage and database directories
# (Using scratch requires specifying permissions elsewhere; adjust if you prefer a non-empty image.)
VOLUME ["/data"]
WORKDIR /data

ENV LISTEN_ADDR=":8080" \
    STORAGE_PATH="/data/storage" \
    DATABASE_URL="/data/files.db" \
    MAX_UPLOAD_SIZE="1073741824"

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/fileserver"]
