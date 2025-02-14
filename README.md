___Запуск приложения___
---
1. **-Убедитесь, что PostgreSQL запущен и доступен.**

2. **-Выполните миграции для создания таблицы tasks.**

3. **-Запустите приложение:**
---
```DATABASE_URL="postgres://user:password@localhost:5432/todo?sslmode=disable" go run cmd/main.go```