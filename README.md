# Брокер сообщений RabbitMQ

## Описание
Демонстрационное приложение TODO-списков с использованием брокера сообщений RabbitMQ, демонстрирующее применение паттернов Clean Architecture и современных практик разработки на Go.

## Технологический стек и концепции
- Разработка REST API на Go с использованием gin-gonic/gin
- Clean Architecture и Dependency Injection
- RabbitMQ в качестве брокера сообщений
- Мониторинг с Prometheus и Grafana
- Логирование с Logrus
- Graceful Shutdown
- Docker для развертывания инфраструктуры

## Требования
- Go (последняя стабильная версия)
- Docker и Docker Compose
- Свободные порты:
  - 15672 (RabbitMQ)
  - 9090 (Prometheus)
  - 3000 (Grafana)
  - 2112 (метрики приложения)

## Быстрый старт

### Запуск инфраструктуры
```bash
docker-compose -f deployments/docker-compose.yml up -d
```

### Запуск приложения
```bash
go run cmd/main.go
```

### Проверка работоспособности

1. Проверка метрик приложения:
```bash
curl http://localhost:2112/metrics
```

2. Проверка запущенных сервисов:
```bash
docker ps
```

3. Доступ к веб-интерфейсам:
- RabbitMQ: http://localhost:15672 (login: guest/guest)
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (login: admin/admin)

## Настройка мониторинга

### Конфигурация Grafana

1. Подключение источника данных:
- Configuration -> Data Sources
- Add data source -> Prometheus
- URL: http://prometheus:9090
- Save & Test

2. Создание дашборда:
- "+" в левом меню -> New Dashboard
- Add visualization
- Data source: Prometheus
- Добавьте следующие метрики:

```promql
# Скорость обработки сообщений
rate(todoapp_rabbitmq_messages_processed_total[5m])

# Время обработки (95-й перцентиль)
histogram_quantile(0.95, rate(todoapp_rabbitmq_message_processing_duration_seconds_bucket[5m]))

# Размер сообщений
histogram_quantile(0.95, rate(todoapp_rabbitmq_message_size_bytes_bucket[5m]))

# Количество ошибок
rate(todoapp_rabbitmq_messages_failed_total[5m])
```

## Структура проекта
```
├── cmd/            # Точка входа приложения
├── deployments/    # Docker конфигурации
├── internal/       # Внутренний код приложения
└── pkg/           # Переиспользуемые пакеты
```

## Дополнительная информация

### Логирование
- Используется logrus для структурированного логирования
- Логи содержат информацию о работе RabbitMQ и обработке сообщений

### Graceful Shutdown
Приложение корректно завершает работу при получении системных сигналов, закрывая все соединения и очереди.

### Метрики
Основные метрики доступны по адресу http://localhost:2112/metrics и включают:
- Количество обработанных сообщений
- Время обработки сообщений
- Размер сообщений
- Количество ошибок

## Решение проблем

### Типичные проблемы

1. Недоступность RabbitMQ:
- Проверьте статус контейнера: `docker ps`
- Проверьте логи: `docker logs rabbitmq`

2. Проблемы с метриками:
- Убедитесь, что Prometheus запущен и доступен
- Проверьте конфигурацию в prometheus.yml

3. Grafana не отображает данные:
- Проверьте подключение к Prometheus
- Убедитесь в корректности PromQL запросов

## Разработка

Для добавления новых функций:
1. Следуйте принципам Clean Architecture
2. Используйте Dependency Injection
3. Добавляйте новые метрики при необходимости
4. Обеспечьте покрытие тестами

---

Для получения дополнительной информации обратитесь к документации используемых технологий:
- [RabbitMQ](https://www.rabbitmq.com/documentation.html)
- [Prometheus](https://prometheus.io/docs/introduction/overview/)
- [Grafana](https://grafana.com/docs/)
- [gin-gonic/gin](https://gin-gonic.com/docs/)