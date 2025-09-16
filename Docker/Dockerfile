# docker/Dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# Copiar todos os arquivos
COPY . .

# Compilar com mais verbosidade
RUN echo "Compilando servidor..." && \
    go build -o luthibox-server main.go && \
    echo "Compilando cliente..." && \
    go build -o luthibox-client client/client.go && \
    echo "Compilação concluída!"

# Verificar se os binários foram criados
RUN ls -la luthibox-* 2>/dev/null || echo "Binários não encontrados"

# Expor porta
EXPOSE 8080

# Comando padr