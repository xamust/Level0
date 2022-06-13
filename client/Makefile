.PHONY: init build

init:
		go mod tidy

build:
		make init
		go build -o build/client -v .
