# Loyalty system (loystem)

## Behaviour schema

```mermaid
sequenceDiagram
autonumber
User->>Loystem: Регистрируется
User-->>Market: Совершает покупку
Order-)CountSystem: Попадает в систему
User->>Loystem: Передаёт номер заказа
Loystem->>Loystem: Связывается номер заказа с пользователем
Loystem->>Loystem: Начисление баллов если есть что
User->>Loystem: Списывает свои баллы за покупки
```

## Run

```shell
./cmd/gophermart -d='postgresql://localhost/postgres?user=postgres&password=postgres'
```

## Links

### Graceful shutdown

* <https://habr.com/ru/articles/771626/>
* <https://followtheprocess.github.io/posts/graceful_shutdown/>
* <https://www.sobyte.net/post/2021-10/go-http-server-shudown-done-right/>
