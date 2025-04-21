DB_URL=postgres://postgres:postgres@0.0.0.0:5432/pvz_service?sslmode=disable
MIGRATE=migrate -path ./migrations -database "$(DB_URL)"

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

migrate-up-n:
	$(MIGRATE) up $(n)

migrate-down-n:
	$(MIGRATE) down $(n)

migrate-force:
	$(MIGRATE) force

migrate-version:
	$(MIGRATE) version

seed:
	docker exec -i postgres psql -U postgres -d pvz_service < ./seeds/seed_test_data.sql

run-tests:
	chmod +x scripts/run-tests.sh
	scripts/run-tests.sh