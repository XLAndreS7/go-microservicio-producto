# Etapa 1: Construcción del binario
FROM golang:1.23 as builder

WORKDIR /app

# Copiar los archivos necesarios
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar el binario
RUN go build -o main ./cmd/main.go

# Etapa 2: Imagen ligera para ejecutar
FROM debian:bookworm-slim

# Crear usuario sin privilegios
RUN useradd -m appuser

# Directorio de trabajo
WORKDIR /app

# Copiar binario desde la etapa anterior
COPY --from=builder /app/main .

# Copiar archivo .env
COPY .env .

# Cambiar dueño del archivo
RUN chown -R appuser:appuser /app

# Cambiar usuario
USER appuser

# Puerto expuesto
EXPOSE 8084

# Comando de ejecución
CMD ["./main"]

