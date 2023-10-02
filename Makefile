.PHONY: build
build:
	cd client && npm install && npm run build
	go build .

.PHONY: build-linux
build-linux:
	cd client && npm install && npm run build
	GOOS=linux GOARCH=amd64 go build .

