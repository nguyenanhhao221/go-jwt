build: 
	go build -o bin/go-jwt

run: build
	./bin/go-jwt

test: 
	go test -v ./...
