# commands

build-cli:
	rm -rf build/ && go build -o build/otter

build-web:
	rm -rf build/web/ && go build -o build/web/otter-web

# start main server
serve:
	go run web/main.go

push:
	git push heroku master

logs:
	heroku logs --tail

# restart dyno
restart:
	heroku ps:restart web
