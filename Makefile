.PHONY: test test-clean

main:
	go build -o build/lang cmd/lang/main.go

test:
	go test ./test/...

m:
	go build -gcflags=-m -o build/lang cmd/lang/main.go

clean:
	rm -f $(TARGET)
	mkdir build

test-clean:
	go clean -testcache
