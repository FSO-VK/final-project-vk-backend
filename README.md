# Бекенд-репозиторий выпускного проекта команды FSO

## Инструменты, которые следует установить

1. [golangci-lint](https://github.com/golangci/golangci-lint) - линтеры

2. [gotestsum](https://github.com/gotestyourself/gotestsum) - форматированный вывод информации о тестах

## Команды

Локальный запуск (для разработки)

```docker compose -f compose.dev.yml```

Форматирование кода

```make format```

Запуск линтеров

```make lint```

Тесты

```make test```
