
.PHONY: all
all: build build/mysql build/wait-for-it
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

build/wait-for-it:
	@go get; go build; mv wait-for-it build/

build/mysql:
	@cd plugins/mysql; go get; go build; mv mysql ../../build/mysql

build:
	@mkdir -p build

