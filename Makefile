all: clean build

build:
	@go build .

watch:
	@echo "Watching for changes. Open up localhost:9999 to see build status instantly."
	@LOG=* go run yolo.go -i *.go -i src -e yolo -c 'make' -a :9999

clean:
	@go clean
