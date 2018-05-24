
.PHONY: all
all: build build/mysql
	@ls -lh build

.PHONY: clean
clean:
	@rm -rf build

.PHONY: copy
copy:
	mkdir -p .wait-for-it/plugins
	cp build/* .wait-for-it/plugins/

.PHONY: run
run:
	go run *.go

build/mysql:
	@cd plugins/mysql; go get; go build; mv mysql ../../build/mysql

build:
	@mkdir -p build

