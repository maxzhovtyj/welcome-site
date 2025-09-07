deploy:
	GOOS=linux GOARCH=amd64 go build -o ./bin/wedding-linux-amd64 ./cmd
	bash release/start.sh