# gophermart
Накопительная система лояльности «Гофермарт»
Индивидуальный дипломный проект курса «Go-разработчик»

# ТЗ
https://practicum.yandex.ru/learn/go-advanced/courses/34ab8873-f88d-4ec4-93cc-c199202c2602/sprints/90962/topics/fb14666e-2249-4b27-8c43-03e17324123a/lessons/fe0e5a74-4022-4cf4-832f-a334e30bafad/


# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без
   префикса `https://`) для создания модуля

## Запуск проекта:

1. Выполнить docker-compose up -d в папке проекта
2. Запустить сервис расчета баллов лояльности ./path/to/project/cmd/accrual/accrual_darwin_arm64 -a localhost:8081
   файл выбирать в зависимости от архитектуры (linux ubuntu - accrual_linux_amd64)
3. В Goland Add Configuration -> go build
4. Run kind = Directory; Directory = к значению, что ide прописало автоматически, надо добавить ```/cmd/gophermart```
5. ENVIRONMENT скопировать из ```.env.server-example```


## Автотесты:
Тесты брать из https://github.com/Yandex-Practicum/go-autotests/releases/tag/v0.10.2
суффикс файла gophermart
для macos - darwin64
для linux - gophermarttest
Пример скрипта запуска автотестов (лучше через Goland configuration). Script Text

Важно!
Перед запуском автотестов, скомпилировать проект в cmd/gophermart/gophermart

```
./gophermarttest
-test.v
-test.run=^TestGophermart$
-gophermart-binary-path=cmd/gophermart/gophermart
-gophermart-host=localhost
-gophermart-port=8080
-gophermart-database-uri="postgres://user:password@localhost:5438/gophermart_postgres?sslmode=disable"
-accrual-binary-path=cmd/accrual/accrual_darwin_arm64
-accrual-host=localhost
-accrual-port=8082
-accrual-database-uri="postgres://user:password@localhost:5438/gophermart_postgres?sslmode=disable"
```

# Линтер
```
go
vet
-vettool=/Users/borisov/GoProjects/yandex/gophermart/statictest
./...
```

# e2e тесты

1. Запустить docker-compose up -d

2. Скомпилировать main.go в cmd/gophermart/gophermart (папку можно настроить свою, но тогда и в перменных окружения запуска теста ее также нужно будет задать)

3. В Goland Add Configuration -> go test

4. Run kind = Package;

5. ENVIRONMENT скопировать из ```.env.e2e-example```
