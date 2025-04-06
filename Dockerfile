FROM golang:1.24-alpine AS builder
LABEL author="masb0ymas"
LABEL name="gofi"

WORKDIR /temp-build

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN cp .env.docker-production .env

# Build the Go app
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o ./cmd/gofi main.go

# Start fresh from a smaller image
FROM alpine:3.21 AS runner
LABEL author="masb0ymas"
LABEL name="gofi"

# Install dependencies for runner
RUN apk add ca-certificates tzdata nano

# Set the Current Working Directory inside the container
WORKDIR /app

# Setup Timezone
ENV TZ=Asia/Jakarta

# Copy built artifacts and necessary files from builder
COPY --from=builder /temp-build/cmd/gofi ./cmd/gofi
COPY --from=builder /temp-build/bin ./bin
COPY --from=builder /temp-build/public ./public
COPY --from=builder /temp-build/secret ./secret
COPY --from=builder /temp-build/.env ./.env

# This container exposes port 8000 to the outside world
EXPOSE 8000

# Run the binary program produced by `go install`
CMD ["./cmd/gofi"]
