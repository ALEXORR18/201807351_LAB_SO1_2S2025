# === Builder Stage ===
#FROM rust:1.90-alpine AS builder

# Instalar dependencias de compilación
#RUN apk add --no-cache protobuf gcc musl-dev make

# Crear directorio de trabajo
#WORKDIR /app

# Copiar archivos de proyecto
#COPY . .

# Compilar la aplicación en modo release
#RUN cargo install --path .

# === Runtime Stage ===
#FROM alpine:3.18

# Crear directorio de trabajo
#WORKDIR /app

# Copiar binario compilado desde el builder
#COPY --from=builder /usr/local/cargo/bin/rust-api ./rust-api

# Exponer el puerto de la API
#EXPOSE 8080

# Variables de entorno
#ENV GRPC_SERVER_URL=http://server-go-service:50051

# Usuario no root
#RUN adduser -D app_user
#USER app_user

# Ejecutar la API
#CMD ["./rust-api"]

# Etapa de build
FROM rust:alpine AS builder

WORKDIR /app

# Instalar dependencias necesarias
RUN apk add --no-cache protobuf gcc musl-dev bash curl

# Copiar los archivos de Rust
COPY . .

# Compilar la app en release
RUN cargo build --release

# Imagen final más ligera
FROM alpine:latest

WORKDIR /app

# Copiar binario compilado desde builder
COPY --from=builder /app/target/release/rust-api ./rust-api

# Exponer puerto de la API
EXPOSE 8080

# Variable de entorno para el server gRPC
ENV GRPC_SERVER_URL=http://go-service:50051

# Crear usuario no root
RUN adduser -D app_user
USER app_user

# Comando de inicio
CMD ["./rust-api"]
