# Initialize the project

keys:
	mkdir -p .keys
	openssl ecparam -genkey -name secp521r1 -noout -out .keys/ecdsa-private.pem
	openssl ec -in .keys/ecdsa-private.pem -pubout -out .keys/ecdsa-public.pem

logdir:
	mkdir -p logs
	touch logs/app.log

configfile:
	cp .gostarter.yaml.example .gostarter.yaml

githook:
	git config core.hooksPath .githooks

init: keys logdir configfile swagger githook
	go mod tidy

clean:
	rm -rf .keys
	rm -rf logs
	rm -rf bin
	rm -rf docs
	rm -rf platform/volumes

install:
	go install github.com/spf13/cobra-cli@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/air-verse/air@latest

# Code Generation

migration:
	go run . migrate create -d platform/migration -n $(name)

migrateup:
	go run main.go migrate up

gqlgen:
	go run github.com/99designs/gqlgen generate

swagger:
	swag init

build:
	GOOS=linux GOARCH=amd64 go build -o bin/gostarter main.go && chmod +x bin/gostarter

# Run

maiden: install init
	docker compose -f docker-compose.dev.yml up --build -d
	go run . serve

dev:
	docker compose -f docker-compose.dev.yml up -d
	air
