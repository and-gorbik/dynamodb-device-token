## dynamodb-device-token

### Используемые фичи
- поддержка нескольких девайсов у одного пользователя (девайсы различаются по device_model)
- TTL на записи в таблице
- Streams для межрегиональной репликации
(- DAX для кеширования)


### Модель данных

device:
- PartitionKey(user_id int64)
- SortKey(kind, device_model string) или SortKey("latest_device")
- modified_at int64
- token string
- app_version string
- locale string
- ttl int64

> Device model нужен в ключе, чтобы поддержать фичу мультидевайс

Проблемы:
- Для однозначного определения токена нужно знать device_model. В текущей реализации мы знаем device model только во время UpdateUserData. Для SendPush придется вытаскивать токены по user_id и kind (при помощи фильтра `begins_with`), после чего 
брать последний по modified_at (на стороне бека).
- При апдейте локали тоже нужно знать device model (сейчас видимо баг)

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

#### Включить стрим
```bash
go run ./cmd/db-manager --command enable-stream
```

возвращает stream id

#### Добавить элемент

```bash
go run ./cmd/api --command put --data '{"user_id": 10, "modified_at": 12345, "kind": "android_general", "device_model": "redmi note 5", "token": "AAA-BBB-CCC-DDDEF", "app_version": "", "locale": "ru"}'
```

#### Добавить элемент и отметить его последним
```bash
go run ./cmd/api --command put --data '{"user_id": 10, "modified_at": 12345, "kind": "android_general", "device_model": "redmi note 5", "token": "AAA-BBB-CCC-DDDEF", "app_version": "", "locale": "ru", "latest": true}'
```

#### Получить несколько записей
```bash
go run ./cmd/api --command get --pk 1 --sort 'ios_general'
```

#### Удалить одну запись
```bash
go run ./cmd/api --command delete --pk 1 --sort 'ios_general#iphone 13'
```

#### Добавить реплику в другой регион
```bash
go run ./cmd/db-manager --command make-global-table --region us-east-1 --replica-region eu-central-1
```

### Запуск в Docker

```bash
docker run -it --rm --name dynamodb-local -p 8000:8000 amazon/dynamodb-local:latest -jar DynamoDBLocal.jar -sharedDb
npm install -g dynamodb-admin
dynamodb-admin
```