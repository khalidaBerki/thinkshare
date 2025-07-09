# Étape 1 : Build de l'application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copie des fichiers go mod et sum
COPY go.mod go.sum ./
RUN go mod download

# Copie du code source
COPY . .

# Compilation de l'application
RUN go build -o thinkshare main.go

# Étape 2 : Image finale
FROM alpine:latest

WORKDIR /app

# Copie du binaire compilé
COPY --from=builder /app/thinkshare .

# Port exposé (à adapter selon ton app)
EXPOSE 8080

# Commande de lancement
CMD ["./thinkshare"]
