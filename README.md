# Серверная часть онлайн чата с использованием технологии Websocket
- Работа с Websocket осуществляется с помощью [gorilla/websocket](https://github.com/gorilla/websocket)
- Пользователи и сообщения сохраняются в PostgreSQL
- Подход чистой архитектуры
- Для авторизации, создания чатов используется [gin/gonic](https://github.com/gin-gonic/gin)
- Представлена базовая работа с Redis, занесение и удаление пользователей по ID чата и пользователя
- Конфигурация PostgreSQL и Redis загружается из .env файла с помощью [joho/godotenv](https://github.com/joho/godotenv)
- Миграция БД с помощью [pressly/goose](https://github.com/pressly/goose)
- PostgreSQL и Redis поднимаются в Docker compose
