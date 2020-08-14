_build_daemon:
	go build -o bin/guardiand cmd/guardiand/main.go

build: _build_daemon

install:
	cp bin/* /usr/local/bin/
