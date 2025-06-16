# Этап сборки фронтенда
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY apps/frontend/ .

# Этап сборки бэкенда
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
COPY . .
RUN apk add --no-cache build-base && \
    cd apps/backend && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/backend ./cmd

# Финальный образ
FROM alpine:3.19
WORKDIR /app

# Устанавливаем ВСЕ необходимые компоненты (включая gcc и g++)
RUN apk add --no-cache \
    nginx \
    tzdata \
    bash \
    curl \
    python3 \
    py3-pip \
    build-base \         
    gcc \              
    g++ && \              
    rm -f /usr/bin/python && \
    ln -s /usr/bin/python3 /usr/bin/python

# Настраиваем Nginx
RUN mkdir -p /run/nginx && \
    mkdir -p /etc/nginx/conf.d

# Копируем файлы
COPY --from=frontend-builder /app /usr/share/nginx/html
COPY --from=backend-builder /go/bin/backend /app/backend
COPY nginx.conf /etc/nginx/nginx.conf

# Создаем симлинки для совместимости
RUN ln -s /labs /app/labs && \
    ln -s /labs /usr/share/nginx/html/labs

EXPOSE 80 8080

CMD ["sh", "-c", "/app/backend -port 8080 & nginx -g 'daemon off;'"]