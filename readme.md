# Messenger

### Требования
- Docker и Docker Compose
- Или: Go 1.23+, PostgreSQL 16

### Развёртывание с Docker Compose

```bash
docker-compose up --build
```

### Развёртывание локально

```bash
go mod download

docker run -e POSTGRES_DB=messenger -e POSTGRES_USER=messenger \
  -e POSTGRES_PASSWORD=messenger -p 5432:5432 postgres:16

export DATABASE_URL="postgres://messenger:messenger@localhost:5432/messenger"

go run cmd/server/main.go
```

### Переменные окружения

```env
DATABASE_URL=postgres://user:password@localhost:5432/messenger
```