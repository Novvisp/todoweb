# Базовый образ с Go
FROM golang:1.23

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Загружаем зависимости
RUN go mod tidy

# Собираем бинарник
RUN go build -o todo-app

# Указываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./todo-app"]
