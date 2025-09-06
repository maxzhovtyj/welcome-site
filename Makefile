deploy:
	go build -i bin/wedding-linux-amd64 cmd/main.go
	bash release/start.sh