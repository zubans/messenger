build:
	docker-compose build

run:
	docker-compose up

stop:
	docker-compose down

migrate-up:
    docker run --rm -v `pwd`/migrations:/liquibase/changelog \
    --network host liquibase/liquibase \
    --url=jdbc:postgresql://localhost:5432/videoconference \
    --changeLogFile=/liquibase/changelog/changelog.xml \
    update

compile-proto:
	docker run --rm --platform=linux/arm64 -v "$(PWD)/proto:/proto" namely/protoc-all \
	  -f /proto/service.proto -l go -o /proto/generated

.PHONY: build run stop migrate-up compile-proto