FROM golang:1.24 AS builder

WORKDIR /app

# Instalar goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copiar archivos y descargar dependencias
COPY db/migrations ./db/migrations
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Compilar binario de la app
RUN go build -o main .

# Imagen final
FROM debian:bookworm-slim

WORKDIR /app

# Instalamos netcat
RUN apt-get update && apt-get install -y netcat-openbsd && rm -rf /var/lib/apt/lists/*

# Copiamos el binario de la app
COPY --from=builder /app/main .

# ⬇️ Copiamos goose desde el builder
COPY --from=builder /go/bin/goose /usr/local/bin/goose
# Copiamos migraciones
COPY --from=builder /app/db/migrations ./db/migrations

COPY start.sh .
RUN chmod +x start.sh

ENTRYPOINT ["./start.sh", "/app/main"]
