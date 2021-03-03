# commands

build:
	rm -rf build/ && go build -o build/otter && cp .env build/.env

# start main server
serve:
	go run web/main.go

push:
	git subtree push --prefix web/ heroku master

logs:
	heroku logs --tail

# restart dyno
restart:
	heroku ps:restart web
