local-test-env-up:
	docker-compose -f ./scripts/docker/docker-compose-local-test.yml up -d
	migrate -source "file://./scripts/db/migrations" -database "postgres://local-test-user:local-test-pass@localhost:5433/critique-local-test?sslmode=disable" up

local-test-env-down:
	docker-compose -f ./scripts/docker/docker-compose-local-test.yml down

local-test:
	make local-test-env-up
	go test ./... -p=1 --cover
