main:
	go build -o build/lang cmd/lang/main.go

clean:
	rm -f $(TARGET)
	mkdir build
