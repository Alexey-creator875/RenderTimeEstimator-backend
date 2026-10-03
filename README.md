# Расчет времени рендера 3D-сцены на рендер-ферме

REST API для публикации и просмотра конфигураций серверов рендеринга.

## HTTP методы

### GET http://localhost:8080/api/render-server-units

Список опубликованных услуг с фильтрацией. Возвращает только записи со статусом `published`. Для каждой записи отдаёт признак `0/1` — создана ли она текущим пользователем.

Query-параметры (фильтрация по объёму оперативной памяти): `min`, `max`.

Ответ `200`:

```json
{
  "renderServerUnits": [
    {
      "id": 1,
      "processor": "Intel Xeon E5-2699",
      "cores": 32,
      "ram": 128,
      "description": "Рендер-ферма под Blender",
      "image": "1_image.png",
      "video": "1_video.mp4",
      "is_mine": 1
    }
  ]
}
```

`is_mine`: `1` — если `creator_id` совпадает с текущим пользователем, иначе `0`.

---

### GET http://localhost:8080/api/render-server-units/:id

Получить одну услугу по ID. Возвращает только опубликованные записи.

Ответ `200`:

```json
{
  "renderServerUnit": {
    "id": 1,
    "processor": "Intel Xeon E5-2699",
    "cores": 32,
    "ram": 128,
    "description": "Рендер-ферма под Blender",
    "image": "1_image.png",
    "video": "1_video.mp4",
    "likes": 42,
    "liked_by_me": true
  }
}
```

Ошибки: `404` — запись не найдена или не опубликована.

---

### GET http://localhost:8080/api/draft-render-server-units

Получить черновик текущего пользователя. ID в запросе не указывается — берётся из аутентификации. У пользователя не может быть больше одного черновика.

Ответ `200`, если черновик есть:

```json
{
  "renderServerUnit": {
    "id": 7,
    "processor": "AMD EPYC 7543",
    "status": "draft"
  }
}
```

Ответ `200`, если черновика нет:

```json
{ "renderServerUnit": null }
```

---

### POST http://localhost:8080/api/render-server-units

Создать черновик. Изображение и короткое видео передаются **как файлы** (multipart/form-data), а не как URL. Имена файлов генерируются на латинице, сами файлы сохраняются в MinIO, а в БД пишутся только имена.

Content-Type: `multipart/form-data`.

Поля формы:

| Поле | Тип | Обязательное | Описание |
|---|---|---|---|
| `processor` | string | да | Модель процессора |
| `image` | file | нет | Изображение (`image/jpeg`, `image/png`, `image/gif`, `image/webp`) |
| `video` | file | нет | Короткое видео (`video/mp4`, `video/webm`, `video/quicktime`) |

Ответ `201`:

```json
{
  "renderServerUnit": {
    "id": 7,
    "processor": "AMD EPYC 7543",
    "status": "draft",
    "image": "7_image.png",
    "video": "7_video.mp4"
  },
  "message": "Черновик успешно добавлен"
}
```

Ошибки: `400` — форма не парсится, файл не того типа, `processor` пустой. `500` — ошибка БД или MinIO.

---

### PUT http://localhost:8080/api/render-server-units

Опубликовать черновик текущего пользователя. Меняет статус на `published` и заполняет оставшиеся поля. ID не передаётся — публикуется черновик текущего пользователя.

Content-Type: `application/json`.

Тело:

```json
{
  "description": "Рендер-ферма под Blender",
  "cores": 32,
  "ram": 128
}
```

Ответ `200`:

```json
{
  "renderServerUnit": {
    "id": 7,
    "status": "published",
    "description": "Рендер-ферма под Blender",
    "cores": 32,
    "ram": 128
  },
  "message": "Юнит опубликован"
}
```

Ошибки: `400` — тело не парсится, `cores` или `ram` ≤ 0. `404` — черновика нет.

---

### DELETE http://localhost:8080/api/render-server-units/:id

Удалить услугу (soft delete). Физически запись остаётся в БД, но помечается удалённой. Удалить можно только собственную запись.

Ответ `200`:

```json
{
  "status": "success",
  "message": "Услуга успешно удалена"
}
```

Ошибки: `403` — попытка удалить чужую запись. `404` — запись не найдена.

---

### POST http://localhost:8080/api/render-server-units/:id/like

Поставить или снять лайк от текущего пользователя.

Content-Type: `application/json`.

Тело:

```json
{ "value": 1 }
```

| Значение | Действие |
|---|---|
| `1` | Поставить лайк (если ещё не стоит) |
| `0` | Снять лайк (если стоит) |

Метод идемпотентен: повторный `value: 1` не создаёт второй лайк, повторный `value: 0` не ломается.

Ответ `200`:

```json
{
  "status": "success",
  "likes": 42,
  "liked": true
}
```

Ошибки: `400` — `value` не `0` и не `1`. `404` — юнит не найден.

---

### POST http://localhost:8080/api/sign-up

Регистрация нового пользователя.

Content-Type: `application/json`.

```json
{ "login": "alice", "password": "secret" }
```

Ответ `201`:

```json
{ "status": "success", "message": "Пользователь создан" }
```

Ошибки: `400` — логин занят, пароль пустой. `500` — ошибка БД.

---

### POST http://localhost:8080/api/sign-in

Аутентификация. Заглушка для 4-й лабораторной.

Content-Type: `application/json`.

```json
{ "login": "alice", "password": "secret" }
```

Ответ `200`:

```json
{ "status": "success", "message": "Вход выполнен" }
```

---

### POST http://localhost:8080/api/sign-out

Деавторизация. Заглушка для 4-й лабораторной.

Ответ `200`:

```json
{ "status": "success", "message": "Выход выполнен" }
```

---

## Таблицы

### render_server_unit

| Поле | Тип | Комментарий |
|---|---|---|
| `id` | bigint | Автоматическое приращение |
| `status` | varchar(10) | `draft` или `published` |
| `processor` | varchar(30) | Модель процессора |
| `cores` | bigint | Количество ядер |
| `ram` | bigint | Объём ОЗУ, ГБ |
| `description` | varchar(300) | Описание |
| `image` | varchar(100) | Имя файла изображения в MinIO |
| `video` | varchar(100) | Имя файла видео в MinIO |
| `created_at` | timestamptz | Дата создания (auto) |
| `creator_id` | bigint | FK → `users.id` |
| `formed_at` | timestamptz | Дата последнего изменения (auto) |

### user

| Поле | Тип | Комментарий |
|---|---|---|
| `id` | bigint | Автоматическое приращение |
| `login` | varchar(40) | Логин |
| `password` | varchar(40) | Пароль |

### like

| Поле | Тип | Комментарий |
|---|---|---|
| `id` | bigint | Автоматическое приращение |
| `user_id` | bigint | FK → `users.id` — кто лайкнул |
| `render_server_unit_id` | bigint | FK → `render_server_units.id` — что лайкнул |

Составной уникальный индекс `(user_id, render_server_unit_id)` — один пользователь может поставить юниту только один лайк.