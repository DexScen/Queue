
# Queue

**Queue** — это веб-приложение для управления очередями и записями пользователей.  
Проект занял 3 место в рамках кейса **ТБанка** на хакатоне **ФакториХак**.

---

## 🚀 Основная идея

Система позволяет пользователям:
- просматривать доступные очереди (например, игры, события, услуги);
- записываться в очередь;
- следить за своим местом;
- администраторам — управлять созданием и ведением очередей.

---

## 🧩 Архитектура проекта

Проект состоит из нескольких сервисов, разворачиваемых через Docker Compose:


Queue/

├── backend/            # Go-сервер (REST API)
│   ├── cmd/main.go     # Точка входа в приложение

│   ├── internal/

│   │   ├── repository/ # Работа с PostgreSQL

│   │   └── transport/  # REST-обработчики и middleware

│   ├── .env.example    # Пример конфигурации окружения

│   ├── Dockerfile

│   └── go.mod, go.sum

│

├── frontend/           # Веб-интерфейс (HTML, CSS, JS)

│   ├── login/          # Страница входа

│   ├── registration/   # Регистрация пользователей

│   ├── user/           # Интерфейс обычного пользователя

│   ├── admin/          # Интерфейс администратора

│   ├── stands/         # Страница со стендами / списками

│   ├── assets/         # Изображения и ресурсы

│   ├── nginx/          # Конфигурация Nginx

│   └── Dockerfile

│

├── PostgreSQL/         # Образ БД и скрипты инициализации

│   ├── Dockerfile

│   └── init.sql

│

├── docker-compose.yml  # Инфраструктура проекта

└── README.md

````

---

## ⚙️ Технологии

**Backend:**

- Go (Golang)

- PostgreSQL (через `database/sql`)

- REST API

- Middleware для авторизации и обработки запросов

**Frontend:**

- HTML + CSS + JavaScript

- Чистая вёрстка без фреймворков

- Адаптивный интерфейс

- Разделение по ролям (пользователь / админ)

**Инфраструктура:**
- Docker, Docker Compose

- Nginx для фронтенда

- Персистентное хранение данных PostgreSQL

---

## 🐳 Запуск проекта

### 1. Установите зависимости

Необходимы:

- [Docker](https://docs.docker.com/get-docker/)

- [Docker Compose](https://docs.docker.com/compose/)

### 2. Склонируйте репозиторий

```bash

git clone https://github.com/DexScen/Queue.git

cd Queue
````

### 3. Создайте `.env` в `backend/` по примеру:

```bash

cp backend/.env.example backend/.env

```

Пример содержимого:

```
DB_HOST=postgres

DB_PORT=5432

DB_USER=postgres

DB_PASSWORD=postgres

DB_NAME=queue_db

PORT=8080
```

### 4. Запустите через Docker Compose

```bash
docker-compose up --build
```

### 5. После успешного запуска

* Frontend будет доступен по адресу:
  👉 [http://localhost:80](http://localhost:80)
  
* Backend API:
  👉 [http://localhost:8080](http://localhost:8080)

---

## 🔗 Основные REST endpoints

| Метод                         | Путь                            | Описание |

| ----------------------------- | ------------------------------- | -------- |

| `GET /api/games`              | Получить список всех очередей   |          |

| `GET /api/games/:id`          | Получить информацию об очереди  |          |

| `POST /api/games/:id/join`    | Присоединиться к очереди        |          |

| `POST /api/auth/register`     | Регистрация нового пользователя |          |

| `POST /api/auth/login`        | Авторизация                     |          |

| `GET /api/users/:login/games` | Список игр пользователя         |          |

(конкретные маршруты можно уточнить в `handler.go`)

---

## 🧱 Структура базы данных

SQL-инициализация задаётся в `PostgreSQL/init.sql` и включает таблицы:

* `users` — информация о пользователях
* `queue` — данные об участниках
* `games` — описания игр/очередей

---

## 💡 Возможности для доработки

* Добавление уведомлений о приближении очереди
* Система ролей (user/admin)
* JWT-авторизация и refresh-токены
* WebSocket-уведомления об обновлении очереди
* Логирование и метрики через Prometheus + Grafana

---

## 👥 Авторы

* **DexScen** — backend, инфраструктура
* **meh-pwn** — frontend
* Команда хакатона **ФакториХак (Т-банк)**
```
