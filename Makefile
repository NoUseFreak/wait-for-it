
build = cd $(1); go get; GOOS=$(2) GOARCH=$(3) go build -ldflags=-s -o build/$(2)_$(4)
build_plugin = cd plugins/$(3); go get; GOOS=$(1) GOARCH=$(2) go build -ldflags=-s -o ../../build/$(1)_$(3)


.PHONY: all
all: test darwin linux
	@ls -lh build

.PHONY: clean
clean:
	@rm -rf build/

.PHONY: copy
copy:
	mkdir -p .wait-for-it/plugins
	cp build/* .wait-for-it/plugins/

.PHONY: test
test:
	go get; \
		go get github.com/stretchr/testify/assert; \
		go test

.PHONY: run
run:
	go run `ls -1 *.go | grep -v _test.go`

.PHONY: install
install:
	mv build/`uname`_wait-for-it /usr/local/bin/wait-for-it
	chmod +x /usr/local/bin/wait-for-it

.PHONY: darwin
darwin: \
	build/darwin_wait-for-it \
	build/darwin_mysql \
	build/darwin_redis \
	build/darwin_cassandra \
	build/darwin_kafka \
	build/darwin_mongodb

.PHONY: linux
darwin: \
	build/linux_wait-for-it \
	build/linux_mysql \
	build/linux_redis \
	build/linux_cassandra \
	build/linux_kafka \
	build/linux_mongodb

build/darwin_wait-for-it:
	$(call build,.,darwin,amd64,wait-for-it)
build/linux_wait-for-it:
	$(call build,.,linux,amd64,wait-for-it)


build/darwin_mysql:
	$(call build_plugin,darwin,amd64,mysql)
build/linux_mysql:
	$(call build_plugin,linux,amd64,mysql)

build/darwin_redis:
	$(call build_plugin,darwin,amd64,redis)
build/linux_redis:
	$(call build_plugin,linux,amd64,redis)

build/darwin_cassandra:
	$(call build_plugin,darwin,amd64,cassandra)
build/linux_cassandra:
	$(call build_plugin,linux,amd64,cassandra)

build/darwin_kafka:
	$(call build_plugin,darwin,amd64,kafka)
build/linux_kafka:
	$(call build_plugin,linux,amd64,kafka)

build/darwin_mongodb:
	$(call build_plugin,darwin,amd64,mongodb)
build/linux_mongodb:
	$(call build_plugin,linux,amd64,mongodb)
