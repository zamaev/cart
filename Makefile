build-all:
	cd cart && GOOS=linux GOARCH=amd64 make build


run-all: build-all
	docker-compose up --force-recreate --build -d

up-monitor:
	sudo docker compose up -d prometheus grafana jaeger

up-db:
	sudo docker compose up -d db_master db_replica

up-infra:
	sudo docker compose up -d db_master db_replica prometheus grafana jaeger kafka-ui kafka0 kafka-init-topics

re: down rm up

up:
	sudo docker compose up -d

down:
	sudo docker compose down

rm:
	sudo docker image rm cart
	sudo docker image rm loms
	sudo docker image rm notifier-image

cnt:
	sudo docker container ls -a

img:
	sudo docker image ls -a
