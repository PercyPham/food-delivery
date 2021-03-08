all:
	@echo "\nAvailable commands:\n"
	@echo "   all         : list all commands\n"
	@echo "   dev         : prepare and run app in development mode\n"
	@echo "   dev_build   : build images before stating dev mode\n"
	@echo "   dev_stop    : stop development mode completely\n"
	@echo "   build       : build executable file in local\n"
	@echo "   build_image : build docker image with default name 'food-delivery'"
	@echo "                 append NAME={{custom-image-name}} to change the image name "
	@echo "\n"

dev:
	@docker-compose -f ./deployments/dev/docker-compose.dev.yaml up -d
	@clear
	@docker logs food-delivery-app -f

dev_build:
	@docker-compose -f ./deployments/dev/docker-compose.dev.yaml up --build -d
	@clear
	@docker logs food-delivery-app -f

dev_stop:
	@docker-compose -f ./deployments/dev/docker-compose.dev.yaml down

build:
	@go build -o main cmd/server/main.go

NAME?=food-delivery
build_image:
	@docker build -t $(NAME) -f build/docker/Dockerfile .

.PHONY: all dev dev_build dev_stop build build_image
