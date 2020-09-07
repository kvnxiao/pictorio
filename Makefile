ifeq ($(OS),Windows_NT)
	OUTPUT = pictorio.exe
else
	OUTPUT = pictorio
endif

build:
	go build -o $(OUTPUT)

lint:
	golangci-lint run
