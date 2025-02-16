# CoinShop

## Требования

Перед запуском убедитесь, что у вас установлены:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)
- [Make](https://www.gnu.org/software/make/)

## Установка и запуск

### 1. Клонирование репозитория
```sh
git clone https://github.com/forsitet/coinshop.git
cd coinshop
```

### 2. Запуск контейнеров
```sh
make up
```
Сервис будет доступен по адресу: [http://localhost:8080](http://localhost:8080)

## Тестирование

Для запуска тестов используйте:

```sh
go test ./...
```

Или для тестирования отдельных пакетов:

```sh
go test ./internal/handlers
```

## Продакшен

Проект развернут на сервере и доступен по адресу:

[http://dbudin.ru:49480](http://dbudin.ru:49480)

## Статический анализ кода

Проект использует **golangci-lint** для автоматической проверки кода на ошибки, потенциальные проблемы и соответствие кодстайлу.

### Результаты нагрузочного тестировния
![image](https://github.com/user-attachments/assets/448271bb-ebfd-43fd-aa04-f6a7dfc1a89e)

