local-test-env-up:
	docker-compose -f ./scripts/docker/docker-compose-local-test.yml up -d

local-test-env-down:
	docker-compose -f ./scripts/docker/docker-compose-local-test.yml down

local-test:
	make local-test-env-up
	go test ./... -count=1
