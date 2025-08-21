# Etapa de compilación
FROM golang:1.22-alpine AS build
WORKDIR /src

# Copiar dependencias
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar el binario
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app

# Imagen final
FROM alpine:3.20
WORKDIR /
COPY --from=build /out/app /app
EXPOSE 3001
ENTRYPOINT ["/app"]
