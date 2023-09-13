SRC = $(shell find . -name \*.go)

.PHONY: default
default: build

.PHONY: build 
build: bin bin/winquit.exe

bin/winquit.exe: export GOOS=windows
bin/winquit.exe: export GOARCH=amd64
bin/winquit.exe: $(SRC) bin
	go build -o bin/winquit.exe ./cmd/winquit

bin:
	mkdir -p bin

.PHONY: clean
clean:
	rm -rf bin

.PHONY: test
test:
	go test -v ./test
