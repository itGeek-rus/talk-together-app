.PHONY: docker-up docker-down docker-logs kafka-topics

#Запуск всех сервисов
docker-up:
    docker-compose up -d --build

#Остановка всех сервисов
docker-down:
    docker-compose down

#Просмотр логов
docker-logs:
    docker-compose logs -f

#Создание топика Kafka (выполнить после запуска)
kafka-topics:
    docker exec kafka-talk kafka-topics --create --topic chat-messages
    --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1

#Остановка и очистка
docker-clean:
    docker-compose down -v
    docker system prune -f

#Проверка статуса сервисов
docker-status:
    docker-compose ps

#Запуск мониторинга
monitoring-up:
    docker-compose -f docker-compose.monitoring.yml up -d

#Остановка мониторинга
monitoring-down:
    docker-compose -f docker-compose.monitoring.yml down

#Полный запуск (приложение + мониторинг)
full-up: docker-up monitoring-up
    @echo "Full stack started!"
    @echo "App: http://localhost:8000"
    @echo "Prometheus: http://localhost:9090"
    @echo "Grafana: http://localhost:3000 (admin/admin)"

#Полная остановка
full-down: docker-down monitoring-down

#Локальная разработка
build:
    go build -o bin/main ./cmd/main.go

run: build
    ./bin/main

dev:
    go run ./cmd/main.go

test:
    go test ./... -v

swagger:
    swag init -g cmd/main.go