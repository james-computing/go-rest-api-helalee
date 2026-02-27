curl -X POST --json '{"title": "todo title", "completed": false}' http://localhost:8080/todos
curl -X POST --json '{"title": "another title", "completed": false}' http://localhost:8080/todos

curl -X GET http://localhost:8080/todos

curl -X GET http://localhost:8080/todos/1

curl -X PUT --json '{"title: "updated title"}' http://localhost:8080/todos/1

curl -X DELETE http://localhost:8080/todos/2