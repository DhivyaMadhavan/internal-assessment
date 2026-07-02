# ---------- Build Stage ----------
FROM golang:1.26.4 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# FIX: Added CGO_ENABLED=0 and fixed output destination folder
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/backend-bin ./backend

# ---------- Runtime Stage ----------
FROM alpine:latest
WORKDIR /app

# FIX: Copy the exact binary generated above
COPY --from=builder /app/backend-bin ./backend

CMD ["./backend"]
