.PHONY: build test clean docker-build

NAME=app
DB_FILE=db.bin

build:
	go build -o ${NAME} github.com/gerlacdt/db-example/cmd/server

test:
	go test -v github.com/gerlacdt/db-example/pkg/...

clean:
	rm -f ${NAME} ${DB_FILE} ./pkg/db/db.test.bin

docker-build:
	GOOS=linux go build -o ${NAME}
	docker build -t gerlacdt/go-example:latest .
	rm -f ${NAME}
