generate-mock:
ifndef FILE
	@echo "FILE parameter is missing"
	@exit 1
endif
	make build-dev
	@docker run --volume "$(PWD)"/../:/app --workdir /app \
	dev-build /bin/bash -c "mockgen -source=${FILE} -destination=mocks/mock_${FILE}"

tests:
	make build-dev
	@docker run \
			--rm \
			--volume "$(PWD)"/../:/app \
			--workdir /app \
			dev-build go test ./... -cover -count=1

swag-fmt:
#	swag fmt -g ../../pkg/server/routes.go -d ../internal/handlers
	make build-dev
	@docker run \
			  --rm \
			  --volume "$(PWD)"/../:/app \
			  --workdir /app \
			  dev-build swag fmt -g ../../pkg/server/routes.go -d ./internal/handlers

swag: swag-fmt
	#swag init -g ../../pkg/server/routes.go -d ../internal/handlers -o ../docs/swagger --parseDependency --parseDepth 1
	make build-dev
	@docker run \
    			--rm \
    			--volume "$(PWD)"/../:/app \
    			--workdir /app \
				dev-build swag init -g ../../pkg/server/routes.go -d ./internal/handlers -o ./docs/swagger --parseDependency --parseDepth 1

start-app:
	docker-compose --file ../deployment/docker-compose.yml down
	docker-compose --file ../deployment/docker-compose.yml up

kill-app:
	docker-compose --file ../deployment/docker-compose.yml down --remove-orphans --volumes

rebuild-app:
	docker-compose --file ../deployment/docker-compose.yml stop app
	docker-compose --file ../deployment/docker-compose.yml build app
	docker-compose --file ../deployment/docker-compose.yml up -d app

build-dev:
	@docker build \
			--tag dev-build \
			-f ../build/Utilities.Dockerfile ..