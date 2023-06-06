# Запуск

Для запуска контейнеров с postgres'ом
```shell
make run_postgres
```

Для запуска контейнеров с redis'ом
```shell
make run_redis
```
Для остановки контейнеров с postgres'ом
```shell
make stop_postgres
```

Для остановки контейнеров с redis'ом
```shell
make stop_redis
```

Для запуска тестового покрытия
```shell
make cover_out
```

Для запуска тестового покрытия и отображения в html
```shell
make cover_html
```

Для генерации protobuf объектов и функций
```shell
make generate_proto
make generate_proto_rpc
```


# API

`docs/swagger.yaml` `docs/swagger.json`

# Порты

|key|default|
|---|---|
|Api Port| 8081|
|Links Port| 8082|
|Postgres Port| 5432|
|Redis Port| 6379|

