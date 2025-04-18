DB_URL=postgres://postgres:postgres@0.0.0.0:5432/pvz_service?sslmode=disable
MIGRATE=migrate -path ./migrations -database "$(DB_URL)"

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

migrate-force:
	$(MIGRATE) force

migrate-version:
	$(MIGRATE) version

