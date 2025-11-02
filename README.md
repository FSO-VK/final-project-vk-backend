# Бекенд-репозиторий выпускного проекта команды FSO

## Инструменты, которые следует установить

1. [golangci-lint](https://github.com/golangci/golangci-lint) - линтеры

2. [gotestsum](https://github.com/gotestyourself/gotestsum) - форматированный вывод информации о тестах

## Команды

Локальный запуск (для разработки). Можно добавить флаг `--watch`
для отслеживания изменений в файлах и пересборки контейнеров.

```docker compose -f compose.dev.yml up```

Форматирование кода

```make format```

Запуск линтеров

```make lint```

Тесты

```make test```
