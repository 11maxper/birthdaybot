# Используем официальный образ Go как базовый
FROM golang:1.25-alpine as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы с зависимостями
COPY go.mod go.sum ./

RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN go build -o bot ./cmd/bot

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

RUN apk add --no-cache tzdata

# Добавляем исполняемый файл из первой стадии в корневую директорию контейнера
COPY --from=builder /app/bot /app/bot

# Указываем, что база и .env будут рядом с бинарником
ENV DB_PATH=/app/data/birthdays.db

# Создаём папку для базы
RUN mkdir -p /app/data

# Запускаем бота
CMD ["./bot"]
