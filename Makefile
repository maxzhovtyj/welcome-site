build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/wedding-linux-amd64 ./cmd

deploy:
	bash release/start.sh