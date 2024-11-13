## dynamodb-device-token

### Примеры cli команд

#### Создать БД

```bash
go run ./cmd/db-manager --command create
```

#### Удалить БД

```bash
go run ./cmd/db-manager --command delete
```

#### Заполнить таблицу значениями

> Таблица заполняется значениями из файла `records.json`.

```bash
go run ./cmd/db-manager --command apply
```

#### Добавить / изменить один элемент

```bash
go run ./cmd/api --command put --data '{"user_id": 10, "modified_at": 12345, "kind": "android_general", "token": "AAA-BBB-CCC-DDDEF", "app_version": "", "device_model": ""}'
```

#### Получить одну запись
```bash
go run ./cmd/api --command get --pk 1 --sort 'android_general'
```

#### Получить несколько записей
```bash
go run ./cmd/api --command get --pk 1
```

#### Удалить одну / несколько записей
```bash
go run ./cmd/api --command delete --pk 1 \[--sort 'ios_general'\]
```