# Этап сборки фронтенда
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY apps/frontend/ .

# Этап сборки бэкенда
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
COPY . .
RUN cd apps/backend && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/backend ./cmd

# Финальный образ
FROM alpine:3.19
WORKDIR /app

# Устанавливаем только необходимые зависимости
RUN apk add --no-cache \
    nginx \
    tzdata \
    build-base \ 
    bash

# Копируем файлы
COPY --from=frontend-builder /app /usr/share/nginx/html
COPY --from=backend-builder /go/bin/backend /app/backend
COPY nginx.conf /etc/nginx/nginx.conf

# Создаем симлинки для совместимости
RUN ln -s /labs /app/labs && \
    ln -s /labs /usr/share/nginx/html/labs

# Проверяем установку компиляторов
RUN gcc --version && g++ --version

EXPOSE 80 8080

CMD ["sh", "-c", "/app/backend -port 8080 & nginx -g 'daemon off;'"]