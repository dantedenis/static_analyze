# Rules for host
host_all: host_build host_run

BINARY_NAME=static_analyze
COVER_FILE=coverage.out

host_build:
	go mod download
	go build -o ${BINARY_NAME} ./cmd/main.go

host_run:
	./${BINARY_NAME}

host_clean:
	go clean
	rm ${BINARY_NAME}

# Testing

lint:
	golangci-lint run

test:
	 go test ./...

coverage:
	go test -coverprofile ${COVER_FILE_TEMP} ./...
	cat ${COVER_FILE_TEMP} | grep -v '.pb.go|mock_' > ${COVER_FILE}
	go tool cover -html=${COVER_FILE}
	rm ${COVER_FILE_TEMP}

test_clean:
	rm ${COVER_FILE}

# Rules for docker

run: build up
kill: stop rm
restart: stop up

create_network:
	docker network create connect-service

build:
	docker-compose -f docker-compose.yml build

up:
	docker-compose -f docker-compose.yml up -d

stop:
	docker-compose -f docker-compose.yml stop

ps:
	docker-compose -f docker-compose.yml ps

rm:
	docker-compose -f docker-compose.yml down

logs:
	docker-compose -f docker-compose.yml logs --tail=100 -f