include .env

compose-up:
	docker-compose up -d --build

compose-down:
	docker-compose down

keys:
	openssl genpkey -algorithm RSA -out private.key && \
	openssl rsa -pubout -in private.key -out public.key
