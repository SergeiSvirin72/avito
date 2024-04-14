devUp:
	docker compose -f ./docker-compose.yaml stop
	docker compose -f ./docker-compose.test.yaml stop
	docker compose -f ./docker-compose.yaml up -d

devMigrationUp:
	docker exec -it ssvirin_app migrate -path=./migration/ -database postgres://postgres:password@ssvirin_db:5432/postgres?sslmode=disable up

testUp:
	docker compose -f ./docker-compose.yaml stop
	docker compose -f ./docker-compose.test.yaml stop
	docker compose -f ./docker-compose.test.yaml up -d

testMigrationUp:
	docker exec -it ssvirin_app_test migrate -path=./migration/ -database postgres://postgres:password@ssvirin_db_test:5432/postgres?sslmode=disable up

testRunTests:
	docker exec -it ssvirin_app_test go test ./cmd/...

allStop:
	docker compose -f ./docker-compose.yaml stop
	docker compose -f ./docker-compose.test.yaml stop