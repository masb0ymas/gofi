# syntax=docker/dockerfile:1

# The build stage runs on the builder's own architecture and cross-compiles,
# so a multi-platform image needs no emulation.
FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/ ./cmd/api ./cmd/migrate

# Static image with CA certificates and tzdata, no shell, running as nonroot.
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=build /out/api /out/migrate ./
# The API serves ./public relative to its working directory; the SQL
# migrations are embedded in the binaries.
COPY public ./public

EXPOSE 8080
ENTRYPOINT ["/app/api"]
