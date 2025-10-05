FROM golang:1.24.5-alpine

WORKDIR /app

#Установка зависимостей
RUN apk add --no-cache gcc musl-dev

#Копирование go mod и sum файлов
COPY go.mod ./

RUN go mod download

#Скопировать исходный код
COPY . .

#Сборка приложения
RUN go build -o main ./cmd/app

#Открыть порт
EXPOSE 8000

#Запуск приложения
CMD ["./main"]