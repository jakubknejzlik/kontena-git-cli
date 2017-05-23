install:
	go get ./...
	go build -o kontena-git

deploy-local:
	make install
	mv kontena-git /usr/local/bin/
