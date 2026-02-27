#!/bin/bash
migrate -path ./migrations -database "postgresql://postgres:replacebypassword@localhost:5432/todo_app_yt?sslmode=disable" up