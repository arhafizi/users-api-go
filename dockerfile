# 1- Build stage
FROM golang:1.24.1-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

RUN wget https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64 -O /usr/local/bin/dbmate \
    && chmod +x /usr/local/bin/dbmate

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

# 2- Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary and config
COPY --from=builder /app/main .
COPY --from=builder /app/config/config-docker.yml ./config/
COPY --from=builder /app/db/migrations ./db/migrations/
COPY --from=builder /usr/local/bin/dbmate /usr/local/bin/dbmate

RUN mkdir -p /app/logs

EXPOSE 5000

CMD ["sh", "-c", "dbmate -d ./db/migrations -e DATABASE_URL up && ./main"]