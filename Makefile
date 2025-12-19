docker_up:
	docker compose -f docker/docker-compose.yaml up

docker_upd:
	docker compose -f docker/docker-compose.yaml up -d

docker_updb: build docker_upd

docker_down:
	docker compose -f docker/docker-compose.yaml down

lint_shortener:
	cd shortener && \
		golint ./... && \
		golangci-lint run ./...

unittest_shortener:
	cd shortener && \
		go test ./tests -v --coverprofile=./tests/cover.out --coverpkg=./pkg/pkgports/adapters/cache/lru && \
		go tool cover --html=./tests/cover.out -o ./tests/cover.html

docker_integration_test:
	cd integration_tests && \
		docker compose up -d && \  # всё кроме e2e_test
		docker compose up --build e2e_test
	cd integration_tests && \
		docker compose down

build:
	docker build -t shortener -f docker/service.Dockerfile ./shortener

init_env:
	type ".\config\example.env" > ".\config\.env"
