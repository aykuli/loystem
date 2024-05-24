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
DATABASE_URI='postgresql://localhost/postgres?user=postgres&password=postgres' ACCRUAL_SYSTEM_ADDRESS='http://localhost:8081' RUN_ADDRESS='localhost:8080' ./cmd/gophermart 

```

## Links

### Graceful shutdown

* <https://github.com/gofiber/fiber/issues/899>
* <https://habr.com/ru/articles/771626/>
* <https://followtheprocess.github.io/posts/graceful_shutdown/>
* <https://www.sobyte.net/post/2021-10/go-http-server-shudown-done-right/>

## Orders Info Poller

```mermaid
graph TD
O[Create poller p] -->A
  A[Возьми все заявки в нужных статусах] --> B
  B[Сделай параллельно запросы\n в систему расчета баллов\n по каждому заказу в количестве p.limit] ---> C
  C{StateCode==429 ?\nToo many requests} ---yes--> D
  C -- no --> A
  D[Descrease p.limit about 1] --> A
```

## TODOs

* FindAll orders rename
* to read
  * time.Tick(p.interval)
  * time.After(p.interval)
