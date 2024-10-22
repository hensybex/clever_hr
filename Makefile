# Makefile

.PHONY: up down restart logs fpm-log ap migration migrate dev

up:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build;

down:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down

restart:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml build api
	docker compose -f docker-compose.yml -f docker-compose.dev.yml restart api
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs -f api

logs:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs -f api

fpm-log:
	docker compose log php-fpm -f

ap:
	docker compose exec php-fpm bash

migration:
	docker compose exec php-fpm php ../bin/console make:migration

cache:
	docker compose exec php-fpm php ../bin/console cache:clear

migrate:
	docker compose exec php-fpm php ../bin/console doctrine:migrations:migrate

dev:
	docker compose -f docker-compose.dev.yml up -d