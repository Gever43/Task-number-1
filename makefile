# Makefile для создания миграций

# Переменные которые будут использоваться в наших командах (Таргетах)
DB_DSN := "postgres://postgres:yourpassword@localhost:5432/main?sslmode=disable" 
MIGRATE := migrate -path ./migrations -database $(DB_DSN)

# Таргет для создания новой миграции для сообщений
migrate-new-messages:
	migrate create -ext sql -dir ./migrations $(NAME)

# Таргет для создания новой миграции для пользователей
migrate-new-users:
	migrate create -ext sql -dir ./migrations $(NAME)

# Применение миграций
migrate:
	$(MIGRATE) up

# Откат миграций
migrate-down:
	$(MIGRATE) down
    
# Для удобства добавим команду run, которая будет запускать наше приложение
run:
	go run cmd/app/main.go

gen_messages:
	oapi-codegen -config openapi/messages/.openapi -include-tags messages -package messages openapi/messages/openapi.yaml > ./internal/web/messages/api.gen.go

lint:
	golangci-lint run --out-format=colored-line-number

gen_users:
	oapi-codegen -config openapi/users/.openapi -include-tags users -package users openapi/users/openapi.yaml > ./internal/web/users/users.go
