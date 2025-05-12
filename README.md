# Распределенный вычислитель арифметических выражений на Golang

Распределенный вычислитель арифметических выражений на Golang 
состоит из 2х основных компонентов: *оркестратор* - "менеджер", 
управляющий и координирующий вычислениями; *агент* - "вычислитель", отвечающий 
за обработку вычислений.
Реализована регистрация и авторизация пользователя. 
При отправке выражений проверяется достоверность пользователя с помощью `JWT токена`. 
Если пользователь проверку не прошел, то он не получит доступ к работе с backtnd-ом.

# Что сделано по задаче
1. Весь реализованный ранее функционал работает как раньше, только в контексте конкретного пользователя.
2. Обеспечено хранение выражений в SQLite. (теперь наша система должна пережить перезагрузку).
3. Общение вычислителя и сервера вычислений реализовано с помощью GRPC.
4. В проекте присутствуют модульные тесты.
5. В проекте присутствуют интеграционные тесты.

## Содержание
### - [Установка](#установка)
### - [Структура проекта](#структура-проекта)
### - [Использование с помощью curl](#использование-с-помощью-curl)

## Установка
Клонировать репозиторий на локальный компьютер
```bash
git clone https://github.com/TravellerGSF/grpc_distr_calc.git
```
Перейти в папку проекта 
```bash
cd grpc_distr_calc
```
1 вариант запуска:

1. Собрать проект командой `docker-compose build`
2. Запустить проект командой `docker-compose up`
3. Для остановки проекта в терминале нажать комбинацию клавиш `Ctrl+C`
4. Перейти на [localhost:8080](http://localhost:8080/) и начать проверять!

2 вариант запуска:

1. Установите обязательные зависимости: `go mod tidy`
2. Запустите оркестратор с рабочей директории: `go run ./cmd/orchestrator/main.go`
3. Запустите агента с рабочей директории: `go run ./cmd/agent/main.go`
4. Перейти на [localhost:8080](http://localhost:8080/) и начать проверять !

## Структура проекта
```
grpc_distr_calc/
├── cmd/ - директория для запуска приложений
│   ├── agent/ - запуск агента (вычислительного узла)
│   └── orchestrator/ - запуск оркестратора (менеджера вычислений)
├── db/ - база данных
│   └── storage.db - SQLite база данных для хранения пользователей и выражений
├── frontend/ - веб-интерфейс
│   ├── auth/ - страницы аутентификации
│   └── main/ - главная страница приложения
├── internal/ - внутренние пакеты проекта
│   ├── grpc/ - gRPC сервисы
│   ├── http/ - HTTP обработчики и middleware
│   ├── storage/ - работа с базой данных
│   ├── test/integration/ - интеграционный тест
│   └── utils/ - вспомогательные утилиты
├── proto/ - протобуфер файлы для gRPC сервисов
├── .env - переменные окружения для конфигурации приложения
├── coverage/ - данные о покрытии тестами
├── docker-compose.yml - конфигурация Docker Compose для развертывания приложения
├── Dockerfile.Agent - Dockerfile для сборки образа агента
├── Dockerfile.Orchestrator - Dockerfile для сборки образа оркестратора
├── go.mod - файл зависимостей Go
├── go.sum - контрольные суммы зависимостей Go
└── README.md - документация проекта
```
## Использование с помощью curl

Регистрация пользователя:
```bash
curl -X POST http://localhost:8080/auth/signup/ \
-H "Content-Type: application/json" \
-d '{"username": "testuser", "password": "testpass"}'
```
Авторизация пользователя:
```bash
curl -X POST http://localhost:8080/auth/login/ \
-H "Content-Type: application/json" \
-d '{"username": "testuser", "password": "testpass"}'
```
После успешной авторизации вы получите JWT токен, который нужно будет использовать для последующих запросов.

Создание выражения:
```bash
curl -X POST http://localhost:8080/expression/ \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
-d '{"expression": "2 + 2"}'
```
Получение всех выражений пользователя:
```bash
curl -X GET http://localhost:8080/expression/ \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```
Ошибочные случаи

Неверные данные при регистрации
```bash
curl -X POST http://localhost:8080/auth/signup/ \
-H "Content-Type: application/json" \
-d '{"username": "testuser", "password": ""}'
```
Неверные данные при авторизации
```bash
curl -X POST http://localhost:8080/auth/login/ \
-H "Content-Type: application/json" \
-d '{"username": "nonexistentuser", "password": "wrongpass"}'
```
Отправка некорректного выражения
```bash
curl -X POST http://localhost:8080/expression/ \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
-d '{"expression": "2 +"}'
```
Эти примеры помогут вам понять, как взаимодействовать с API вашего распределённого калькулятора.





