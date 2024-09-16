# Используем официальный образ Golang
FROM golang:1.23.0

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем Air
RUN go install github.com/air-verse/air@latest

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем код приложения
COPY . .

# Собираем бинарный файл
RUN go build -o tender-service .

# Указываем команду для запуска приложения
CMD ["air"]
