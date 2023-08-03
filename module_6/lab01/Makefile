clear:
	rm -rf bin

build:
	go build -o bin/main main.go

run:
	./bin/main --backends=https://grafana.com:443,http://info.cern.ch:80   

hello:
	echo "Hello"

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/main-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go

all: clear build run