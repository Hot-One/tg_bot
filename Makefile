run:
	go run cmd/main.go

migration-up:
	migrate -path ./migration/postgres -database 'postgres://abdulbosit:946236953@localhost:5432/bot?sslmode=disable' up

migration-down:
	migrate -path ./migration/postgres -database 'postgres://abdulbosit:946236953@localhost:5432/bot?sslmode=disable' down


