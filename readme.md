# readme

## Генерация ручек

```shell
go generate -v ./...
```

## Создание миграции
```shell
$migration_name='migration_name'
goose -dir=./internal/migrations create $migration_name go
```
