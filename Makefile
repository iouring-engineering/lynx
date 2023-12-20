PKG_NAME=lynx-services
PKG_VERSION=0.0.1

BIN=build/lynx

all: BIN

install:
	go mod download

BIN: src/main.go
	go build -o ${BIN} src/*.go
	
run:
	go run -race src/*.go -config config/config.yaml

clean:
	rm -f ${BIN}
