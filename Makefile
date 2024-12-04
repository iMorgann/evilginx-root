TARGET=evilginx
PACKAGES=core database log parser

.PHONY: all
all: build run

dev: build debug

build:
	@go mod tidy
	@go mod vendor
	@go build -o ./bin/$(TARGET) -buildvcs=false

clean:
	@go clean
	@rm -f ./bin/$(TARGET)

install:
	@mkdir -p /usr/share/evilginx/phishlets
	@mkdir -p /usr/share/evilginx/templates
	@cp ./phishlets/* /usr/share/evilginx/phishlets/
	@cp ./templates/* /usr/share/evilginx/templates/
	@cp ./bin/$(TARGET) /usr/local/bin

run:
	@sudo ./bin/evilginx -p ./phishlets -c ./config

debug:
	@sudo ./bin/evilginx -developer -p ./phishlets -c ./config -debug
