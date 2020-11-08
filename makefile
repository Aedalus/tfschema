.ONESHELL:
.SHELLFLAGS = -ec

build:
	go build -o bin/main main.go

compile:
	@rm -rf bin

	GOOS=darwin GOARCH=amd64 go build -o bin/tfschema-darwin-amd64 main.go;
	GOOS=darwin GOARCH=386 go build -o bin/tfschema-darwin-386 main.go;

	GOOS=linux GOARCH=amd64 go build -o bin/tfschema-linux-amd64 main.go;
	GOOS=linux GOARCH=386 go build -o bin/tfschema-linux-386 main.go;

	GOOS=windows GOARCH=amd64 go build -o bin/tfschema-windows-amd64 main.go;
	GOOS=windows GOARCH=386 go build -o bin/tfschema-windows-386 main.go;


run:
	go run main.go $(ARGS)