build-all:
	cd cart && GOOS=linux GOARCH=amd64 make build


run-all: build-all
	docker-compose up --force-recreate --build -d

up:
	sudo docker compose up -d

down:
	sudo docker compose down

rm:
	sudo docker image rm cart
	sudo docker image rm loms

cnt:
	sudo docker container ls -a

img:
	sudo docker image ls -a
