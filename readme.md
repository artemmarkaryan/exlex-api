# readme

## SSH

```shell
ssh root@178.21.10.214
```

## Требуемые переменные окружения

- `PSQL_DSN`
- `UNIONE_TOKEN`
- `JWT_SECRET_KEY`

## Генерация ручек

```shell
go generate -v ./...
```

## Миграции

### Установка `goose` (macOS)

```shell
brew install goose
```

### Создание миграции
```shell
$migration_name='migration_name'
goose -dir=./internal/migrations create $migration_name go
```

