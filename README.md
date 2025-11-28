

# Создать пользователя
curl -X POST http://localhost:8080/api/users ^
  -H "Content-Type: application/json" ^
  -d "{\"name\":\"Alice\",\"email\":\"alice@example.com\"}"

# Получить всех
curl http://localhost:8080/api/users

# Получить по id
curl http://localhost:8080/api/users/1

# Обновить
curl -X PUT http://localhost:8080/api/users/1 ^
  -H "Content-Type: application/json" ^
  -d "{\"name\":\"Alice Updated\",\"email\":\"alice2@example.com\"}"

# Удалить
curl -X DELETE http://localhost:8080/api/users/1