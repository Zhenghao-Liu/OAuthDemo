export GO111MODULE=on
BIN_NAME="OAuth_demo"

all:clean build run

build:
	go mod tidy
	bash -x build.sh

run:
	./output/bin/${BIN_NAME}

log:
	less -N output/log/gin.log

clean:
	rm -rf output

test:
	go test -v