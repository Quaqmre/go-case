.PHONY: clean-images build-images build debug run stop dep test

clean-images:
	@echo "---------------- Cleaning dangling Docker images ----------------"
	docker images -f "dangling=true" -q | xargs --no-run-if-empty docker rmi -f

build-images:
	@echo "------------------------ Building Images ------------------------"
	docker compose -f docker-compose.yml build

build: build-images clean-images

debug:
	docker compose -f docker-compose.yml up

run:
	docker compose -f docker-compose.yml up -d

stop:
	docker compose -f docker-compose.yml down

dep:
	go mod download

test:
	ginkgo ./...