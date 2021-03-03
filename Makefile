# commands

build:
	rm -rf build/ && go build -o build/otter

# start main server
serve:
	go run web/main.go