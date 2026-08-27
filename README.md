## Технологии

* Go 1.25
* PostgreSQL
* Telebot v3
* Cron v3
* Docker
* Git
* Linux
* GitHub
* Godotenv

## Развертывание

Проект контейнеризирован с использованием Docker, что позволяет быстро развернуть приложение в изолированной среде и обеспечить одинаковое окружение для разработки и запуска.

### Запуск через Docker

```bash
sudo docker build -t reminder-bot .
sudo docker run --env-file .env reminder-bot
```

### Старт на прод через Docker Compose
#### Скрытый старт
``` bash
sudo docker compose up -d
```
#### Старт в окне терминала
```bash
sudo docker compose up
```

## Ключевые технические особенности

* Telegram Bot API
* PostgreSQL для хранения данных
* Многопоточная обработка задач через goroutines
* Фоновые воркеры для отправки напоминаний
* Cron-задачи для повторяющихся уведомлений
* Конфигурация через переменные окружения
* Docker-контейнеризация
* FSM (Finite State Machine) для управления пользовательскими сценариями

## Roadmap

- CI/CD через GitHub Actions
- Автоматизированное тестирование
- Мониторинг и логирование
- Повышение отказоустойчивости сервиса
