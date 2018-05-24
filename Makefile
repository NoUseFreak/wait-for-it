
build = cd $(1); go get; GOOS=$(2) GOARCH=$(3) go build -ldflags=-s -o build/$(2)_$(4)
build_plugin = cd plugins/$(3); go get; GOOS=$(1) GOARCH=$(2) go build -ldflags=-s -o ../../build/$(1)_$(3)


.PHONY: all
all: darwin linux
	@ls -lh build

.PHONY: clean
clean:
	@rm -rf build/

.PHONY: copy
copy:
	mkdir -p .wait-for-it/plugins
	cp build/* .wait-for-it/plugins/

.PHONY: run
run:
	go run *.go

.PHONY: darwin
darwin: build/darwin_wait-for-it build/darwin_mysql

.PHONY: linux
darwin: build/linux_wait-for-it build/linux_mysql

build/darwin_wait-for-it:
	$(call build,.,darwin,amd64,wait-for-it)

build/linux_wait-for-it:
	$(call build,.,linux,amd64,wait-for-it)


build/darwin_mysql:
	$(call build_plugin,darwin,amd64,mysql)

build/linux_mysql:
	$(call build_plugin,linux,amd64,mysql)

