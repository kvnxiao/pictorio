ifeq ($(OS),Windows_NT)
	OUTPUT = pictorio.exe
else
	OUTPUT = pictorio
endif

build:
	go build -tags production -o $(OUTPUT)

dev:
	go build -tags development -o $(OUTPUT)

lint:
	golangci-lint run
