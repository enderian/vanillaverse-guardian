build:
	go build -ldflags "-X 'main.baseFolder=/etc/guardiand'" -o bin/guardiand cmd/guardiand/main.go
	go build -o bin/guardian cmd/guardian/main.go

install:
	cp bin/* /usr/local/bin/

run:
	go run cmd/guardiand/main.go

cli:
	go run cmd/guardian/main.go
