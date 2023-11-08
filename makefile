include .env
export


build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

generate-certs:
	openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout ./certs/cert.key -out ./certs/cert.crt

build-image:
	docker build -t gptutor-bot .

start-container:
	docker run --env-file .env -p 80:80 -p 443:443 -e TZ=Europe/Moscow -v /GPTutorBot-config/db:/GPTutorBot/db -v /GPTutorBot-config/certs:/GPTutorBot/certs gptutor-bot

up:
	docker-compose up -d
