install:
	go get ./...
	go build -o kontena-git
	mv kontena-git /usr/local/bin/
