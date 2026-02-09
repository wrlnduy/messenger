# Messenger

### Требования
- Docker и Docker Compose

### Развёртывание с Docker Compose

```bash
docker-compose up --build
```

### Переменные окружения

```env
DATABASE_URL=postgres://user:password@postgres:5432/messenger
REDIS_URL=redis://redis:6379/0
```