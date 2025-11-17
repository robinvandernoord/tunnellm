# Makefile
.SILENT:

build:
	go build -o target/tunellm .

run: build
	./target/tunellm
