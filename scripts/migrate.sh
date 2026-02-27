#!/bin/bash
# usage: from scripts folder, run
# ./migrate.sh up
# or
# ./migrate.sh down
migrate -path ../migrations -database "postgresql://postgres:replacebypassword@localhost:5432/todo_app_yt?sslmode=disable" $1