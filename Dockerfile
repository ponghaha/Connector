# -----------------------------
# Stage 1: Build
# -----------------------------
FROM golang:1.24.12-alpine AS builder

ENV GOPROXY=https://proxy.golang.org,direct

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# -----------------------------
# Stage 2: Runtime
# -----------------------------
FROM alpine:3.20

RUN apk add --no-cache tzdata curl
ENV TZ=Asia/Bangkok

WORKDIR /root

COPY --from=builder /app/cmd/server/server .
COPY configs/ /root/configs/
COPY migration/ /root/migration/
COPY master/ /root/master/

EXPOSE 8080
CMD ["./server"]
