BINARY_NAME=line-oa-manager

build:
	go build -o bin/${BINARY_NAME} cmd/http/main.go
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin cmd/http/main.go
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux cmd/http/main.go
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}-windows cmd/http/main.go
	GOARCH=arm GOOS=linux go build -o bin/${BINARY_NAME}-arm cmd/http/main.go
	GOARCH=arm64 GOOS=linux go build -o bin/${BINARY_NAME}-arm64 cmd/http/main.go

run:
	./${BINARY_NAME}

clean:
	go clean
	rm bin/*