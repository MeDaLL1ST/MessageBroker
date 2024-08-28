# Используйте официальный образ Golang как базовый
FROM golang:1.22-alpine as builder

# Скопируйте исходный код в рабочую директорию внутри контейнера
WORKDIR /app

# Копируйте go.mod и go.sum в рабочую директорию
COPY go.mod go.sum ./

# Загрузите зависимости
RUN go mod tidy

# Копируйте остальной исходный код в рабочую директорию
COPY . .

# Соберите приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Используйте alpine образ для упрощения размера конечного образа
FROM alpine:latest  
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Копируйте исполняемый файл из предыдущего этапа
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
EXPOSE 8070
EXPOSE 5137
# Команда запуска приложения
CMD ["./main"]
