.PHONY: build test run proto clean docker-build

NAME=app
DB_FILE=app.db.bin
PB_DIR=./pb

build:
	go build -o ${NAME} github.com/gerlacdt/db-example/cmd/server

test:
	go test github.com/gerlacdt/db-example/...

run: build
	PORT=8080 DB_FILENAME=${DB_FILE} ./app

proto:
	protoc -I ${PB_DIR} ${PB_DIR}/db.proto --go_out=${PB_DIR}

clean:
	rm -f ${NAME} ${DB_FILE} ./pkg/db/db.test.bin

docker-build:
	GOOS=linux go build -o ${NAME} github.com/gerlacdt/db-example/cmd/server
	docker build -t gerlacdt/db-example:latest .
	rm -f ${NAME}
