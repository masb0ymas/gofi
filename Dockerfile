FROM golang:1.24-alpine AS builder
LABEL author="masb0ymas"
LABEL name="gofi"

WORKDIR /temp-build

RUN apk update && apk add --no-cache make

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN cp .envrc.example .envrc

# Build the application
RUN make build/api
RUN make build/migrate

# Create the final image
FROM alpine:3.21

WORKDIR /app
RUN apk add ca-certificates tzdata nano

COPY --from=builder /temp-build/bin/api /app/api
COPY --from=builder /temp-build/bin/migrate /app/migrate
COPY --from=builder /temp-build/templates /app/templates
COPY --from=builder /temp-build/public /app/public
COPY --from=builder /temp-build/.envrc /app/.envrc

EXPOSE 8080
CMD ["./api"]
