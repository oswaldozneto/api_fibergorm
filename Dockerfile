# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instala dependências de build
RUN apk add --no-cache git

# Copia go.mod e go.sum
COPY go.mod go.sum ./

# Download das dependências
RUN go mod download

# Copia o código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Instala certificados CA
RUN apk --no-cache add ca-certificates tzdata

# Copia o binário compilado
COPY --from=builder /app/main .

# Expõe a porta
EXPOSE 3000

# Comando de execução
CMD ["./main"]

