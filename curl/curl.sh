# Public endpoints
curl -X POST --json '{"email": "john@something.com", "password": "password123"}' http://localhost:8080/auth/register

curl -X POST --json '{"email": "john@something.com", "password": "password123"}' http://localhost:8080/auth/login

# Authenticated endpoints

curl -X GET --json '{"email": "john@something.com", "password": "password123"}' -H "Authorization: Bearer replaceByAccessToken" http://localhost:8080/protected-test

curl -X POST --json '{"title": "todo title", "completed": false}' -H "Authorization: Bearer replaceByAccessToken" http://localhost:8080/todos
curl -X POST --json '{"title": "another title", "completed": false}' -H "Authorization: Bearer replaceByAccessToken" http://localhost:8080/todos

curl -X GET -H "Authorization: Bearer replaceByAccessToken" http://localhost:8080/todos

curl -X GET -H "Authorization: Bearer replaceByAccessToken" http://localhost:8080/todos/1

curl -X PUT --json '{"title": "updated title"}' -H "Authorization: Bearer replaceByAccessToken" http://localhost:8080/todos/1

curl -X DELETE -H "Authorization: Bearer replaceByAccessToken" http://localhost:8080/todos/2
