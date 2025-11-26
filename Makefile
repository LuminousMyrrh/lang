main:
	go build -o dist/lang cmd/lang/main.go

clean:
	rm -f $(TARGET)
	mkdir dist
