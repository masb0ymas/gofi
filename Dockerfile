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

RUN cp .envrc.docker .envrc

# Build the application
RUN make build/api
RUN make build/migrate

# Create the final image
FROM alpine:3.21
LABEL author="masb0ymas"
LABEL name="gofi"

WORKDIR /app
RUN apk add ca-certificates tzdata nano make

# Copy the built application
COPY --from=builder /temp-build/bin/api /app/api
COPY --from=builder /temp-build/bin/migrate /app/migrate
COPY --from=builder /temp-build/migrations /app/migrations
COPY --from=builder /temp-build/templates /app/templates
COPY --from=builder /temp-build/public /app/public
COPY --from=builder /temp-build/script/Makefile /app/Makefile
COPY --from=builder /temp-build/.envrc /app/.envrc

# Process .envrc and create .env file during build
RUN touch /app/.env && \
    while IFS= read -r line; do \
    if echo "$line" | grep -q "^export"; then \
    echo "$line" | sed 's/^export //' >> /app/.env; \
    fi \
    done < /app/.envrc && \
    cat /app/.env

# Expose the port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/bin/sh", "-c", "set -a && . /app/.env && set +a && exec ./api \
    --machine-id=$MACHINE_ID \
    --debug=$DEBUG \
    --env=$ENV \
    --port=$PORT \
    --app-name=$APP_NAME \
    --app-secret=$APP_SECRET \
    --jwt-secret=$JWT_SECRET \
    --client-url=$CLIENT_URL \
    --server-url=$SERVER_URL \
    --db-dsn=$DB_DSN \
    --db-max-open-conns=$DB_MAX_OPEN_CONNS \
    --db-max-idle-conns=$DB_MAX_IDLE_CONNS \
    --db-max-idle-time=$DB_MAX_IDLE_TIME \
    --redis-addr=$REDIS_ADDR \
    --redis-password=$REDIS_PASSWORD \
    --redis-db=$REDIS_DB \
    --resend-api-key=$RESEND_API_KEY \
    --resend-from-email=$RESEND_FROM_EMAIL \
    --resend-debug-to-email=$RESEND_DEBUG_TO_EMAIL \
    --google-client-id=$GOOGLE_CLIENT_ID \
    --google-client-secret=$GOOGLE_CLIENT_SECRET \
    --google-redirect-url=$GOOGLE_REDIRECT_URL"]
